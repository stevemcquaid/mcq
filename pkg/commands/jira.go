package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

// Config holds Jira connection configuration
type Config struct {
	URL      string
	Username string
	Password string
}

// Issue represents a Jira issue for display
type Issue struct {
	Key         string
	Summary     string
	Description string
	Status      string
	Assignee    string
	Reporter    string
	Priority    string
	Type        string
	Sprint      string
	Parent      string
	Labels      []string
	Components  []string
	FixVersions []string
	Created     time.Time
	Updated     time.Time
	Comments    []Comment
}

// Comment represents a Jira comment
type Comment struct {
	Author  string
	Body    string
	Created time.Time
	Updated time.Time
}

// UserCache holds a simple cache for resolved usernames
var userCache = make(map[string]string)

// askForConfirmation asks the user for confirmation with a default value
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

// getConfig retrieves Jira configuration from viper
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

// createClient creates a Jira client with basic authentication
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

// ShowJiraIssue displays detailed information about a Jira issue
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

// stripHTML removes basic HTML tags from text and converts Jira links to markdown or plain text
func stripHTML(text string) string {
	// Convert Jira links to markdown or plain text format
	text = convertJiraLinks(text)

	// Convert Jira code blocks and formatting
	text = convertJiraCodeBlocks(text)

	// Simple HTML tag removal - this is basic but covers most cases
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "<p>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "<strong>", "")
	text = strings.ReplaceAll(text, "</strong>", "")
	text = strings.ReplaceAll(text, "<em>", "")
	text = strings.ReplaceAll(text, "</em>", "")
	text = strings.ReplaceAll(text, "<b>", "")
	text = strings.ReplaceAll(text, "</b>", "")
	text = strings.ReplaceAll(text, "<i>", "")
	text = strings.ReplaceAll(text, "</i>", "")

	// Remove any remaining HTML tags (basic regex-like approach)
	// This is a simple implementation - for production use, consider a proper HTML parser
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// Clean up extra whitespace
	text = strings.TrimSpace(text)
	return text
}

// convertJiraLinks converts Jira link format to markdown or plain text
func convertJiraLinks(text string) string {
	// Convert Jira user links [~accountid:...] to @username
	text = convertJiraUserLinks(text)

	// Convert Jira smart links [text|url|smart-link] to markdown [text](url) or just url
	text = convertJiraSmartLinks(text)

	return text
}

// convertJiraCodeBlocks converts Jira code formatting to markdown
func convertJiraCodeBlocks(text string) string {
	// Convert {noformat} blocks to triple backticks
	text = convertJiraNoFormatBlocks(text)

	// Convert {code} blocks to triple backticks
	text = convertJiraCodeBlocksWithLang(text)

	// Convert inline code formatting
	text = convertJiraInlineCode(text)

	// Convert headings
	text = convertJiraHeadings(text)

	return text
}

// convertJiraNoFormatBlocks converts {noformat}...{noformat} to ```...```
func convertJiraNoFormatBlocks(text string) string {
	start := 0
	for {
		// Look for {noformat}
		startTag := strings.Index(text[start:], "{noformat}")
		if startTag == -1 {
			break
		}
		startTag += start

		// Look for closing {noformat}
		endTag := strings.Index(text[startTag+10:], "{noformat}")
		if endTag == -1 {
			break
		}
		endTag += startTag + 10

		// Extract the content between tags
		content := text[startTag+10 : endTag]

		// Replace with markdown code block
		markdownBlock := fmt.Sprintf("```\n%s\n```", content)
		text = text[:startTag] + markdownBlock + text[endTag+10:]

		// Update start position
		start = startTag + len(markdownBlock)
	}

	return text
}

// convertJiraCodeBlocksWithLang converts {code:lang}...{code} to ```lang...```
func convertJiraCodeBlocksWithLang(text string) string {
	start := 0
	for {
		// Look for {code:lang} or {code}
		startTag := strings.Index(text[start:], "{code")
		if startTag == -1 {
			break
		}
		startTag += start

		// Find the end of the opening tag
		tagEnd := strings.Index(text[startTag:], "}")
		if tagEnd == -1 {
			break
		}
		tagEnd += startTag + 1

		// Extract language if present
		lang := ""
		if text[startTag+5:tagEnd-1] != "" {
			lang = text[startTag+6 : tagEnd-1] // Skip ":"
		}

		// Look for closing {code}
		endTag := strings.Index(text[tagEnd:], "{code}")
		if endTag == -1 {
			break
		}
		endTag += tagEnd

		// Extract the content between tags
		content := text[tagEnd:endTag]

		// Replace with markdown code block
		markdownBlock := fmt.Sprintf("```%s\n%s\n```", lang, content)
		text = text[:startTag] + markdownBlock + text[endTag+6:]

		// Update start position
		start = startTag + len(markdownBlock)
	}

	return text
}

// convertJiraInlineCode converts {{code}} to `code`
func convertJiraInlineCode(text string) string {
	start := 0
	for {
		// Look for {{code}}
		startTag := strings.Index(text[start:], "{{")
		if startTag == -1 {
			break
		}
		startTag += start

		// Look for closing }}
		endTag := strings.Index(text[startTag+2:], "}}")
		if endTag == -1 {
			break
		}
		endTag += startTag + 2

		// Extract the content between tags
		content := text[startTag+2 : endTag]

		// Replace with markdown inline code
		markdownCode := fmt.Sprintf("`%s`", content)
		text = text[:startTag] + markdownCode + text[endTag+2:]

		// Update start position
		start = startTag + len(markdownCode)
	}

	return text
}

// convertJiraHeadings converts Jira headings to markdown
func convertJiraHeadings(text string) string {
	// Convert h1. to #
	text = strings.ReplaceAll(text, "h1. ", "# ")

	// Convert h2. to ##
	text = strings.ReplaceAll(text, "h2. ", "## ")

	// Convert h3. to ###
	text = strings.ReplaceAll(text, "h3. ", "### ")

	// Convert h4. to ####
	text = strings.ReplaceAll(text, "h4. ", "#### ")

	// Convert h5. to #####
	text = strings.ReplaceAll(text, "h5. ", "##### ")

	// Convert h6. to ######
	text = strings.ReplaceAll(text, "h6. ", "###### ")

	return text
}

// convertJiraUserLinks converts Jira user account IDs to @username format
func convertJiraUserLinks(text string) string {
	// Find all occurrences of [~accountid:...]
	start := 0
	for {
		// Look for the start of a user link
		userStart := strings.Index(text[start:], "[~accountid:")
		if userStart == -1 {
			break
		}
		userStart += start

		// Find the end of the user link
		userEnd := strings.Index(text[userStart:], "]")
		if userEnd == -1 {
			break
		}
		userEnd += userStart + 1

		// Extract the account ID
		accountID := text[userStart+12 : userEnd-1] // Skip "[~accountid:" and "]"

		// Resolve the account ID to actual username
		username := resolveAccountID(accountID)
		if username == "" {
			// Fallback to a simplified format if resolution fails
			username = "@user-" + accountID[len(accountID)-8:]
		}

		// Replace the user link with @username
		text = text[:userStart] + username + text[userEnd:]

		// Update start position to continue searching
		start = userStart + len(username)
	}

	return text
}

// convertJiraSmartLinks converts Jira smart links to markdown or plain text format
func convertJiraSmartLinks(text string) string {
	// Pattern: [text|url|smart-link] or [text|url] -> [text](url) or just url
	start := 0
	for {
		// Look for the start of a smart link
		linkStart := strings.Index(text[start:], "[")
		if linkStart == -1 {
			break
		}
		linkStart += start

		// Find the end of the link
		linkEnd := strings.Index(text[linkStart:], "]")
		if linkEnd == -1 {
			break
		}
		linkEnd += linkStart + 1

		// Extract the link content
		linkContent := text[linkStart+1 : linkEnd-1] // Skip "[" and "]"

		// Check if this is a smart link format [text|url|smart-link] or [text|url]
		parts := strings.Split(linkContent, "|")
		if (len(parts) == 3 && parts[2] == "smart-link") || (len(parts) == 2) {
			linkText := parts[0]
			linkURL := parts[1]

			// Parse and clean the URL to make it valid
			var decodedURL string
			parsedURL, err := url.Parse(linkURL)
			if err != nil {
				// If parsing fails, use the original URL
				decodedURL = linkURL
			} else {
				// The parsed URL should be properly formatted
				decodedURL = parsedURL.String()
			}

			var convertedLink string
			if linkText == linkURL {
				// If text and URL are the same, just print the decoded URL once
				convertedLink = decodedURL
			} else {
				// If text and URL are different, use markdown format with decoded URL
				convertedLink = fmt.Sprintf("[%s](%s)", linkText, decodedURL)
			}

			text = text[:linkStart] + convertedLink + text[linkEnd:]
			start = linkStart + len(convertedLink)
		} else {
			// Not a smart link, continue searching
			start = linkStart + 1
		}
	}

	return text
}

// resolveAccountID attempts to resolve a Jira account ID to a username
func resolveAccountID(accountID string) string {
	// Check cache first
	if username, exists := userCache[accountID]; exists {
		return username
	}

	// Get configuration for API call
	config, err := getConfig()
	if err != nil {
		// If we can't get config, return empty string (will use fallback)
		return ""
	}

	// Make API call to resolve account ID
	username := fetchUsernameByAccountID(config, accountID)
	if username != "" {
		// Cache the result
		userCache[accountID] = username
		return username
	}

	// Return empty string if resolution failed (will use fallback)
	return ""
}

// fetchUsernameByAccountID makes an API call to resolve account ID to username
func fetchUsernameByAccountID(config *Config, accountID string) string {
	// Construct the API URL
	apiURL := fmt.Sprintf("%s/rest/api/2/user?accountId=%s", config.URL, accountID)

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ""
	}

	// Add authentication
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the error but don't fail the operation
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	// Parse response
	var userResponse struct {
		DisplayName  string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
		Key          string `json:"key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return ""
	}

	// Return the display name, or email if display name is empty
	if userResponse.DisplayName != "" {
		return "@" + userResponse.DisplayName
	} else if userResponse.EmailAddress != "" {
		// Extract username from email (part before @)
		emailParts := strings.Split(userResponse.EmailAddress, "@")
		if len(emailParts) > 0 {
			return "@" + emailParts[0]
		}
	}

	return ""
}

// displayIssue displays the issue information in a formatted way
func displayIssue(issue *Issue) {
	fmt.Printf("\n🔍 Jira Issue: %s\n", issue.Key)
	fmt.Println(strings.Repeat("=", 50))

	// Basic information
	fmt.Printf("📋 Summary: %s\n", stripHTML(issue.Summary))
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
		cleanDescription := stripHTML(issue.Description)
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
				cleanBody := stripHTML(comment.Body)
				cleanBody = strings.ReplaceAll(cleanBody, "\n", "\n   ")
				fmt.Printf("   %s\n\n", cleanBody)
			}
		} else {
			fmt.Println("Skipping comments.")
		}
	}

	fmt.Println(strings.Repeat("=", 50))
}
