package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

// Client provides a simple interface for JIRA API operations
type Client struct {
	client *jira.Client
	config *Config
}

// NewClient creates a new JiraClient instance
func NewClient() (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	client, err := createClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

// GetIssue retrieves a JIRA issue by key
func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	jiraIssue, _, err := c.client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	issue := c.convertJiraIssue(jiraIssue)

	// Fetch comments
	comments, err := c.GetComments(issueKey)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not fetch comments: %v\n", err)
	} else {
		issue.Comments = comments
	}

	return issue, nil
}

// GetComments retrieves comments for a JIRA issue
func (c *Client) GetComments(issueKey string) ([]Comment, error) {
	apiURL := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", c.config.URL, issueKey)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.config.Username, c.config.Password)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

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

	return c.convertComments(commentResponse.Comments), nil
}

// CreateIssue creates a new JIRA issue
func (c *Client) CreateIssue(issue *jira.Issue) (string, error) {
	createdIssue, _, err := c.client.Issue.Create(issue)
	if err != nil {
		return "", fmt.Errorf("failed to create issue: %w", err)
	}

	return createdIssue.Key, nil
}

// GetBaseURL returns the base URL for the Jira instance
func (c *Client) GetBaseURL() string {
	return c.config.URL
}

// convertJiraIssue converts a JIRA issue to our Issue struct
func (c *Client) convertJiraIssue(jiraIssue *jira.Issue) *Issue {
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

	return issue
}

// convertComments converts JIRA comments to our Comment struct
func (c *Client) convertComments(jiraComments []struct {
	Author struct {
		DisplayName string `json:"displayName"`
	} `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}) []Comment {
	var comments []Comment
	for _, c := range jiraComments {
		created, _ := time.Parse("2006-01-02T15:04:05.000-0700", c.Created)
		updated, _ := time.Parse("2006-01-02T15:04:05.000-0700", c.Updated)

		comments = append(comments, Comment{
			Author:  c.Author.DisplayName,
			Body:    c.Body,
			Created: created,
			Updated: updated,
		})
	}
	return comments
}

// getConfig retrieves Jira configuration from viper and environment variables
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

// createClient creates a Jira client with basic authentication using the provided config
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
