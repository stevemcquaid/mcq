// Package commands provides Jira integration functionality for the mcq CLI tool.
// This file contains the core Jira API client, issue fetching, and display logic.
// Text formatting and conversion logic is handled in jira_formatter.go.
package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

// Config holds Jira connection configuration
type Config struct {
	URL      string // Jira instance URL
	Username string // Username for authentication
	Password string // Password or API token for authentication
}

// Issue represents a Jira issue for display with all relevant fields
type Issue struct {
	Key         string    // Issue key (e.g., "PROJ-123")
	Summary     string    // Issue title/summary
	Description string    // Issue description
	Status      string    // Current status
	Assignee    string    // Assigned user
	Reporter    string    // User who created the issue
	Priority    string    // Issue priority
	Type        string    // Issue type (Bug, Story, etc.)
	Sprint      string    // Sprint name if assigned
	Parent      string    // Parent issue key for subtasks
	Labels      []string  // Issue labels
	Components  []string  // Project components
	FixVersions []string  // Fix versions
	Created     time.Time // Creation timestamp
	Updated     time.Time // Last update timestamp
	Comments    []Comment // Issue comments
}

// Comment represents a Jira comment with metadata
type Comment struct {
	Author  string    // Comment author
	Body    string    // Comment content
	Created time.Time // Comment creation time
	Updated time.Time // Comment last update time
}

// Global text formatter instance
var textFormatter = NewTextFormatter()

// askForConfirmation prompts the user for confirmation with a default value.
// If defaultNo is true, the default is "no" (user can press Enter for no).
// If defaultNo is false, the default is "yes" (user can press Enter for yes).
func askForConfirmation(prompt string, defaultNo bool) bool {
	reader := bufio.NewReader(os.Stdin)

	defaultText := "y/N"
	if !defaultNo {
		defaultText = "Y/n"
	}

	fmt.Printf("%s [%s]: ", prompt, defaultText)

	response, err := reader.ReadString('\n')
	if err != nil {
		return !defaultNo
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" {
		return !defaultNo
	}

	return response == "y" || response == "yes"
}

// getConfig retrieves Jira configuration from viper and environment variables.
// Returns an error if required configuration is missing.
func getConfig() (*Config, error) {
	url := viper.GetString("jira.url")
	username := viper.GetString("jira.username")
	password := viper.GetString("jira.password")
	token := viper.GetString("jira.token")

	if url == "" {
		return nil, fmt.Errorf("jira URL not configured. Set JIRA_INSTANCE_URL environment variable or use --url flag")
	}

	if username == "" {
		return nil, fmt.Errorf("jira username not configured. Set JIRA_USERNAME environment variable or use --username flag")
	}

	// Use API token as password if provided
	if token != "" {
		password = token
	}

	if password == "" {
		return nil, fmt.Errorf("jira password/token not configured. Set JIRA_PASSWORD or JIRA_API_TOKEN environment variable")
	}

	return &Config{
		URL:      url,
		Username: username,
		Password: password,
	}, nil
}

// createClient creates a Jira client with basic authentication using the provided config.
func createClient(config *Config) (*jira.Client, error) {
	transport := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}

	client, err := jira.NewClient(transport.Client(), config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return client, nil
}

// normalizeIssueKey adds project prefix if the issue key is just a number.
// If JIRA_PROJECT_PREFIX is set, "123" becomes "PROJ-123".
// If no prefix is set or the key already contains a dash, returns the key as-is.
func normalizeIssueKey(issueKey string) string {
	// If it already contains a dash, assume it's a full key
	if strings.Contains(issueKey, "-") {
		return issueKey
	}

	// Check if it's just a number
	prefix := viper.GetString("jira.project_prefix")
	if prefix == "" {
		// No prefix configured, return as-is
		return issueKey
	}

	// Add the prefix
	return fmt.Sprintf("%s-%s", prefix, issueKey)
}

// ShowJiraIssue displays detailed information about a Jira issue.
// This is the main entry point for the "mcq jira show" command.
func ShowJiraIssue(issueKey string) {
	// Normalize issue key (add prefix if needed)
	normalizedKey := normalizeIssueKey(issueKey)

	config, err := getConfig()
	if err != nil {
		fmt.Printf("❌ Configuration error: %v\n", err)
		os.Exit(1)
	}

	client, err := createClient(config)
	if err != nil {
		fmt.Printf("❌ Failed to create Jira client: %v\n", err)
		os.Exit(1)
	}

	issue, err := fetchIssue(client, normalizedKey)
	if err != nil {
		fmt.Printf("❌ Failed to fetch issue %s: %v\n", normalizedKey, err)
		os.Exit(1)
	}

	displayIssue(issue)
}

// fetchIssue retrieves issue details from Jira
func fetchIssue(client *jira.Client, issueKey string) (*Issue, error) {
	jiraIssue, _, err := client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	issue := &Issue{
		Key:         jiraIssue.Key,
		Summary:     jiraIssue.Fields.Summary,
		Description: jiraIssue.Fields.Description,
		Status:      jiraIssue.Fields.Status.Name,
		Priority:    jiraIssue.Fields.Priority.Name,
		Type:        jiraIssue.Fields.Type.Name,
		Created:     time.Time(jiraIssue.Fields.Created),
		Updated:     time.Time(jiraIssue.Fields.Updated),
	}

	// Optional fields
	if jiraIssue.Fields.Assignee != nil {
		issue.Assignee = jiraIssue.Fields.Assignee.DisplayName
	}
	if jiraIssue.Fields.Reporter != nil {
		issue.Reporter = jiraIssue.Fields.Reporter.DisplayName
	}
	if jiraIssue.Fields.Parent != nil {
		issue.Parent = jiraIssue.Fields.Parent.Key
	}

	// Collections
	issue.Labels = jiraIssue.Fields.Labels
	issue.Components = extractComponentNames(jiraIssue.Fields.Components)
	issue.FixVersions = extractFixVersionNames(jiraIssue.Fields.FixVersions)
	issue.Sprint = extractSprintName(jiraIssue.Fields.Unknowns)

	// Fetch comments
	comments, err := fetchComments(client, issueKey)
	if err != nil {
		// Don't fail the entire operation if comments can't be fetched
		fmt.Printf("⚠️  Warning: Could not fetch comments: %v\n", err)
	} else {
		issue.Comments = comments
	}

	return issue, nil
}

// extractComponentNames extracts component names from Jira components
func extractComponentNames(components []*jira.Component) []string {
	var names []string
	for _, component := range components {
		names = append(names, component.Name)
	}
	return names
}

// extractFixVersionNames extracts version names from Jira fix versions
func extractFixVersionNames(versions []*jira.FixVersion) []string {
	var names []string
	for _, version := range versions {
		names = append(names, version.Name)
	}
	return names
}

// extractSprintName extracts sprint name from custom fields
func extractSprintName(unknowns map[string]interface{}) string {
	sprintField := unknowns["customfield_10020"]
	if sprintField == nil {
		return ""
	}

	sprints, ok := sprintField.([]interface{})
	if !ok || len(sprints) == 0 {
		return ""
	}

	sprint, ok := sprints[0].(map[string]interface{})
	if !ok {
		return ""
	}

	name, exists := sprint["name"]
	if !exists {
		return ""
	}

	sprintName, ok := name.(string)
	if !ok {
		return ""
	}

	return sprintName
}

// fetchComments retrieves comments for an issue using direct API call
func fetchComments(client *jira.Client, issueKey string) ([]Comment, error) {
	// Get the base URL from the client
	baseURL := client.GetBaseURL()
	apiURL := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", baseURL.String(), issueKey)

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers - we need to get these from the config
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the error but don't fail the operation
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var commentResponse struct {
		Comments []struct {
			Author struct {
				DisplayName string `json:"displayName"`
			} `json:"author"`
			Body    string `json:"body"`
			Created string `json:"created"`
			Updated string `json:"updated"`
		} `json:"comments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&commentResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our Comment struct
	var comments []Comment
	for _, c := range commentResponse.Comments {
		created, _ := time.Parse("2006-01-02T15:04:05.000-0700", c.Created)
		updated, _ := time.Parse("2006-01-02T15:04:05.000-0700", c.Updated)

		comments = append(comments, Comment{
			Author:  c.Author.DisplayName,
			Body:    c.Body,
			Created: created,
			Updated: updated,
		})
	}

	return comments, nil
}

// formatText applies all Jira text formatting using the global formatter
func formatText(text string) string {
	return textFormatter.FormatText(text)
}

// All text formatting functions have been moved to jira_formatter.go

// displayIssue displays the issue information in a formatted way
func displayIssue(issue *Issue) {
	fmt.Printf("\n🔍 Jira Issue: %s\n", issue.Key)
	fmt.Println(strings.Repeat("=", 50))

	// Basic information
	fmt.Printf("📋 Summary: %s\n", formatText(issue.Summary))
	fmt.Printf("📝 Type: %s\n", issue.Type)
	fmt.Printf("📊 Status: %s\n", issue.Status)
	fmt.Printf("⚡ Priority: %s\n", issue.Priority)

	// People
	if issue.Assignee != "" {
		fmt.Printf("👤 Assignee: %s\n", issue.Assignee)
	}
	if issue.Reporter != "" {
		fmt.Printf("📢 Reporter: %s\n", issue.Reporter)
	}

	// Dates
	fmt.Printf("📅 Created: %s\n", issue.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("🔄 Updated: %s\n", issue.Updated.Format("2006-01-02 15:04:05"))

	// Optional fields
	if issue.Sprint != "" {
		fmt.Printf("🏃 Sprint: %s\n", issue.Sprint)
	}
	if issue.Parent != "" {
		fmt.Printf("👨‍👩‍👧‍👦 Parent: %s\n", issue.Parent)
	}

	// Collections
	if len(issue.Labels) > 0 {
		fmt.Printf("🏷️  Labels: %s\n", strings.Join(issue.Labels, ", "))
	}
	if len(issue.Components) > 0 {
		fmt.Printf("🧩 Components: %s\n", strings.Join(issue.Components, ", "))
	}
	if len(issue.FixVersions) > 0 {
		fmt.Printf("🔧 Fix Versions: %s\n", strings.Join(issue.FixVersions, ", "))
	}

	// Description
	if issue.Description != "" {
		fmt.Printf("\n📄 Description:\n")
		// Apply the same formatting as comments (HTML cleaning and link conversion)
		cleanDescription := formatText(issue.Description)
		cleanDescription = strings.ReplaceAll(cleanDescription, "\n", "\n")
		fmt.Printf("%s\n", cleanDescription)
	}

	// Comments
	if len(issue.Comments) > 0 {
		fmt.Printf("\n💬 Comments (%d) available.\n", len(issue.Comments))

		// Ask user if they want to see comments
		if askForConfirmation("Show comments?", true) {
			fmt.Println(strings.Repeat("-", 30))
			for i, comment := range issue.Comments {
				fmt.Printf("%d. %s (%s):\n", i+1, comment.Author, comment.Created.Format("2006-01-02 15:04:05"))
				// Clean up the comment body (remove HTML tags and format nicely)
				cleanBody := formatText(comment.Body)
				cleanBody = strings.ReplaceAll(cleanBody, "\n", "\n   ")
				fmt.Printf("   %s\n\n", cleanBody)
			}
		} else {
			fmt.Println("Skipping comments.")
		}
	}

	fmt.Println(strings.Repeat("=", 50))
}

// ============================================================================
// JIRA NEW COMMAND FUNCTIONS
// ============================================================================

// JiraNew creates a new Jira issue from a vague user story using AI
func JiraNew(args []string, modelFlag string, verbosityLevel int, contextConfig ContextConfig) error {
	featureRequest := strings.Join(args, " ")

	// First, generate the user story using the existing AI functionality
	userStory, err := generateUserStoryForJira(featureRequest, modelFlag, verbosityLevel, contextConfig)
	if err != nil {
		return fmt.Errorf("failed to generate user story: %w", err)
	}

	// Ask for confirmation before creating the Jira issue
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Generated User Story:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(userStory)
	fmt.Println(strings.Repeat("=", 60))

	if !askForConfirmation("\nCreate Jira issue with this content?", false) {
		fmt.Println("Jira issue creation cancelled.")
		return nil
	}

	// Create the Jira issue
	issueKey, err := createJiraIssue(userStory, featureRequest)
	if err != nil {
		return fmt.Errorf("failed to create Jira issue: %w", err)
	}

	// Display success message
	fmt.Printf("\n✅ Jira issue created successfully: %s\n", issueKey)
	fmt.Printf("🔗 You can view it at: %s/browse/%s\n", getJiraBaseURL(), issueKey)

	return nil
}

// generateUserStoryForJira generates a user story using AI and returns it without copying to clipboard
func generateUserStoryForJira(featureRequest string, modelFlag string, verbosityLevel int, contextConfig ContextConfig) (string, error) {
	// Set up logger
	SetupLogger(verbosityLevel)
	logBasic("Starting JiraNew", "feature_request", featureRequest)

	// Gather repository context
	repoContext := GatherContextIfNeeded(contextConfig)

	// Select and configure model
	selectedModel, err := SelectModel(modelFlag)
	if err != nil {
		return "", err
	}

	// Generate user story
	userStory, err := GenerateUserStory(selectedModel, featureRequest, repoContext)
	if err != nil {
		return "", fmt.Errorf("failed to generate user story: %w", err)
	}

	// Copy to clipboard (as requested)
	if err := copyToClipboard(userStory); err != nil {
		logError("clipboard copy", err)
		// Don't fail the entire operation if clipboard copy fails
		fmt.Printf("⚠️  Warning: Could not copy to clipboard: %v\n", err)
	} else {
		fmt.Println("📋 User story copied to clipboard!")
	}

	return userStory, nil
}

// createJiraIssue creates a new Jira issue with the provided content
func createJiraIssue(userStory, featureRequest string) (string, error) {
	config, err := getConfig()
	if err != nil {
		return "", fmt.Errorf("configuration error: %w", err)
	}

	// Get project prefix from environment
	projectPrefix := viper.GetString("jira.project_prefix")
	if projectPrefix == "" {
		return "", fmt.Errorf("JIRA_PROJECT_PREFIX environment variable is required")
	}

	client, err := createClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create Jira client: %w", err)
	}

	// Extract title from the user story (first line or first sentence)
	title := extractTitleFromUserStory(userStory, featureRequest)

	// Convert markdown-style formatting to Jira markup
	formattedDescription := convertToJiraMarkup(userStory)

	// Create the issue
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Project:     jira.Project{Key: projectPrefix},
			Type:        jira.IssueType{Name: "Story"}, // Default to Story
			Summary:     title,
			Description: formattedDescription,
		},
	}

	createdIssue, _, err := client.Issue.Create(issue)
	if err != nil {
		return "", fmt.Errorf("failed to create issue: %w", err)
	}

	return createdIssue.Key, nil
}

// extractTitleFromUserStory extracts a suitable title from the user story
func extractTitleFromUserStory(userStory, featureRequest string) string {
	// Try to find the main user story line (starts with "As a")
	lines := strings.Split(userStory, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "As a") {
			// Extract the goal part from "As a [user], I want [goal] so that [benefit]"
			parts := strings.Split(line, "I want ")
			if len(parts) > 1 {
				goalPart := strings.Split(parts[1], " so that")[0]
				return strings.TrimSpace(goalPart)
			}
		}
	}

	// Fallback to the original feature request, cleaned up
	title := strings.TrimSpace(featureRequest)
	if len(title) > 100 {
		title = title[:97] + "..."
	}
	return title
}

// getJiraBaseURL returns the base URL for the Jira instance
func getJiraBaseURL() string {
	config, err := getConfig()
	if err != nil {
		return "https://your-jira-instance.com" // fallback
	}
	return config.URL
}

// convertToJiraMarkup converts markdown-style text to Jira markup
func convertToJiraMarkup(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			result = append(result, "")
			continue
		}

		// Convert bullet points
		if strings.HasPrefix(line, "- ") {
			// Main bullet point
			content := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			result = append(result, "* "+content)
		} else if strings.HasPrefix(line, "  - ") {
			// Sub bullet point (indented)
			content := strings.TrimSpace(strings.TrimPrefix(line, "  - "))
			result = append(result, "** "+content)
		} else if strings.HasPrefix(line, "    - ") {
			// Sub-sub bullet point (double indented)
			content := strings.TrimSpace(strings.TrimPrefix(line, "    - "))
			result = append(result, "*** "+content)
		} else if strings.HasPrefix(line, "1. ") {
			// Numbered list
			content := strings.TrimSpace(strings.TrimPrefix(line, "1. "))
			result = append(result, "# "+content)
		} else if strings.HasPrefix(line, "  1. ") {
			// Indented numbered list
			content := strings.TrimSpace(strings.TrimPrefix(line, "  1. "))
			result = append(result, "## "+content)
		} else if strings.HasPrefix(line, "## ") {
			// H2 heading
			content := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			result = append(result, "h2. "+content)
		} else if strings.HasPrefix(line, "### ") {
			// H3 heading
			content := strings.TrimSpace(strings.TrimPrefix(line, "### "))
			result = append(result, "h3. "+content)
		} else if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") {
			// Bold text
			content := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, "**"), "**"))
			result = append(result, "*"+content+"*")
		} else if strings.HasPrefix(line, "*") && strings.HasSuffix(line, "*") && !strings.HasPrefix(line, "**") {
			// Italic text
			content := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, "*"), "*"))
			result = append(result, "_"+content+"_")
		} else if strings.HasPrefix(line, "```") {
			// Code block start/end
			if strings.HasPrefix(line, "```") && len(line) == 3 {
				result = append(result, "{code}")
			} else if strings.HasPrefix(line, "```") && len(line) > 3 {
				// Code block with language
				lang := strings.TrimSpace(strings.TrimPrefix(line, "```"))
				result = append(result, "{code:"+lang+"}")
			}
		} else if strings.HasPrefix(line, "`") && strings.HasSuffix(line, "`") {
			// Inline code
			content := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, "`"), "`"))
			result = append(result, "{{"+content+"}}")
		} else {
			// Regular text
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
