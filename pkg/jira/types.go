package jira

import (
	"fmt"
	"time"
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

// ValidationError represents a simple validation error
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", ve.Field, ve.Message)
}
