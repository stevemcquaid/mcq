package jira

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

// Manager provides a simplified interface for JIRA operations
type Manager struct {
	client      *Client
	aiExtractor *AIExtractor
}

// NewManager creates a new Manager instance
func NewManager() (*Manager, error) {
	client, err := NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return &Manager{
		client:      client,
		aiExtractor: nil, // Will be set by SetAIExtractor
	}, nil
}

// SetAIExtractor sets the AI extractor for the manager
func (m *Manager) SetAIExtractor(extractor *AIExtractor) {
	m.aiExtractor = extractor
}

// GetBaseURL returns the base URL for the Jira instance
func (m *Manager) GetBaseURL() string {
	return m.client.GetBaseURL()
}

// ShowIssue displays a JIRA issue with all details
func (m *Manager) ShowIssue(issueKey string) error {
	normalizedKey := normalizeIssueKey(issueKey)
	issue, err := m.client.GetIssue(normalizedKey)
	if err != nil {
		return fmt.Errorf("failed to fetch issue %s: %w", normalizedKey, err)
	}

	m.displayIssue(issue)
	return nil
}

// CreateIssue creates a new JIRA issue from a user story
func (m *Manager) CreateIssue(userStory, featureRequest string) (string, error) {
	// Validate inputs
	if err := m.validateInputs(userStory, featureRequest); err != nil {
		return "", err
	}

	// Get project configuration
	projectKey, err := m.getProjectKey()
	if err != nil {
		return "", err
	}

	// Extract title using AI with user approval
	title, err := m.extractTitle(userStory, featureRequest)
	if err != nil {
		return "", fmt.Errorf("failed to extract title: %w", err)
	}

	// Create the issue
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Project:     jira.Project{Key: projectKey},
			Type:        jira.IssueType{Name: "Story"},
			Summary:     title,
			Description: convertToJiraMarkup(userStory),
		},
	}

	issueKey, err := m.client.CreateIssue(issue)
	if err != nil {
		return "", fmt.Errorf("failed to create issue: %w", err)
	}

	return issueKey, nil
}

// extractTitle extracts a title using AI with user approval
func (m *Manager) extractTitle(userStory, featureRequest string) (string, error) {
	// Try AI extraction first
	aiTitle, err := m.extractTitleWithAI(userStory, featureRequest)
	if err != nil {
		fmt.Printf("⚠️  Warning: AI title extraction failed: %v\n", err)
		fmt.Println("Falling back to pattern-based extraction...")
		return extractTitleWithPatterns(userStory, featureRequest), nil
	}

	// Show AI title and get user approval
	fmt.Printf("\n🤖 AI-generated title: \"%s\"\n", aiTitle)

	if askForConfirmation("Use this title for the Jira issue?", false) {
		return aiTitle, nil
	}

	// User rejected, ask for custom title
	customTitle := m.getCustomTitle()
	if customTitle != "" {
		return customTitle, nil
	}

	// Fall back to pattern-based extraction
	return extractTitleWithPatterns(userStory, featureRequest), nil
}

// extractTitleWithAI uses AI to extract a title
func (m *Manager) extractTitleWithAI(userStory, featureRequest string) (string, error) {
	if m.aiExtractor == nil {
		return "", fmt.Errorf("AI functionality not available - AI extractor not set")
	}
	return m.aiExtractor.ExtractTitleWithAI(userStory, featureRequest)
}

// getCustomTitle prompts user for a custom title
func (m *Manager) getCustomTitle() string {
	fmt.Print("Enter custom title (or press Enter to use pattern-based extraction): ")
	reader := bufio.NewReader(os.Stdin)
	customTitle, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to read input: %v\n", err)
		return ""
	}

	customTitle = strings.TrimSpace(customTitle)
	if len(customTitle) > 100 {
		customTitle = customTitle[:97] + "..."
	}

	return customTitle
}

// validateInputs validates the input parameters
func (m *Manager) validateInputs(userStory, featureRequest string) error {
	if strings.TrimSpace(userStory) == "" {
		return ValidationError{Field: "userStory", Message: "cannot be empty"}
	}
	if strings.TrimSpace(featureRequest) == "" {
		return ValidationError{Field: "featureRequest", Message: "cannot be empty"}
	}
	return nil
}

// getProjectKey gets the JIRA project key from configuration
func (m *Manager) getProjectKey() (string, error) {
	projectKey := viper.GetString("jira.project_prefix")
	if projectKey == "" {
		return "", fmt.Errorf("JIRA_PROJECT_PREFIX environment variable is required")
	}
	return projectKey, nil
}

// displayIssue displays issue information in a clean format
func (m *Manager) displayIssue(issue *Issue) {
	formatter := NewTextFormatter()

	fmt.Printf("\n🔍 Jira Issue: %s\n", issue.Key)
	fmt.Println(strings.Repeat("=", 50))

	// Basic info
	fmt.Printf("📋 Summary: %s\n", formatter.FormatText(issue.Summary))
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
	m.displayCollections(issue)

	// Description
	m.displayDescription(issue, formatter)

	// Comments
	m.displayComments(issue, formatter)

	fmt.Println(strings.Repeat("=", 50))
}

// displayCollections displays issue collections (labels, components, etc.)
func (m *Manager) displayCollections(issue *Issue) {
	if len(issue.Labels) > 0 {
		fmt.Printf("🏷️  Labels: %s\n", strings.Join(issue.Labels, ", "))
	}
	if len(issue.Components) > 0 {
		fmt.Printf("🧩 Components: %s\n", strings.Join(issue.Components, ", "))
	}
	if len(issue.FixVersions) > 0 {
		fmt.Printf("🔧 Fix Versions: %s\n", strings.Join(issue.FixVersions, ", "))
	}
}

// displayDescription displays the issue description
func (m *Manager) displayDescription(issue *Issue, formatter *TextFormatter) {
	if issue.Description == "" {
		return
	}

	fmt.Printf("\n📄 Description:\n")
	cleanDescription := formatter.FormatText(issue.Description)
	cleanDescription = strings.ReplaceAll(cleanDescription, "\n", "\n")
	fmt.Printf("%s\n", cleanDescription)
}

// displayComments displays issue comments
func (m *Manager) displayComments(issue *Issue, formatter *TextFormatter) {
	if len(issue.Comments) == 0 {
		return
	}

	fmt.Printf("\n💬 Comments (%d) available.\n", len(issue.Comments))

	if !askForConfirmation("Show comments?", true) {
		fmt.Println("Skipping comments.")
		return
	}

	fmt.Println(strings.Repeat("-", 30))
	for i, comment := range issue.Comments {
		fmt.Printf("%d. %s (%s):\n", i+1, comment.Author, comment.Created.Format("2006-01-02 15:04:05"))
		cleanBody := formatter.FormatText(comment.Body)
		cleanBody = strings.ReplaceAll(cleanBody, "\n", "\n   ")
		fmt.Printf("   %s\n\n", cleanBody)
	}
}

// askForConfirmation prompts the user for confirmation with a default value
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

// normalizeIssueKey adds project prefix if the issue key is just a number
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

// extractTitleWithPatterns uses pattern-based approach to extract title
func extractTitleWithPatterns(userStory, featureRequest string) string {
	// Try to find the main user story line (starts with "As a")
	lines := strings.Split(userStory, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "As a") {
			// Extract the goal part from "As a [user], I want [goal] so that [benefit]"
			parts := strings.Split(line, "I want ")
			if len(parts) > 1 {
				goalPart := strings.Split(parts[1], " so that")[0]
				title := strings.TrimSpace(goalPart)
				if title != "" {
					return title
				}
			}
		}
	}

	// Try alternative patterns
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for lines that start with "I want" or similar patterns
		if strings.HasPrefix(line, "I want") {
			// Extract everything after "I want"
			title := strings.TrimSpace(strings.TrimPrefix(line, "I want"))
			if title != "" {
				return title
			}
		}
		// Look for lines that start with "User should" or similar
		if strings.HasPrefix(line, "User should") {
			title := strings.TrimSpace(line)
			if title != "" {
				return title
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
