package errors

import (
	"fmt"
	"strings"
)

// UserError represents a user-friendly error with guidance
type UserError struct {
	Code         string
	Message      string
	Suggestion   string
	Troubleshoot []string
	OriginalErr  error
}

func (e *UserError) Error() string {
	return e.Message
}

// Unwrap returns the original error for error chaining
func (e *UserError) Unwrap() error {
	return e.OriginalErr
}

// Display shows a formatted error message with troubleshooting steps
func (e *UserError) Display() {
	fmt.Printf("%s\n", e.Message)
	if e.Suggestion != "" {
		fmt.Printf("💡 %s\n", e.Suggestion)
	}
	if len(e.Troubleshoot) > 0 {
		fmt.Println("\n🔧 Troubleshooting steps:")
		for i, step := range e.Troubleshoot {
			fmt.Printf("   %d. %s\n", i+1, step)
		}
	}
}

// Predefined error types
var (
	JiraAuthError = &UserError{
		Code:       "JIRA_AUTH_FAILED",
		Message:    "❌ JIRA Authentication Failed",
		Suggestion: "Check your credentials and try again",
		Troubleshoot: []string{
			"Verify JIRA_INSTANCE_URL is correct",
			"Check JIRA_API_TOKEN is valid",
			"Run 'mcq config test' to validate setup",
			"Run 'mcq config setup' for interactive configuration",
		},
	}

	JiraConfigError = &UserError{
		Code:       "JIRA_CONFIG_MISSING",
		Message:    "❌ JIRA Configuration Missing",
		Suggestion: "Set up JIRA configuration first",
		Troubleshoot: []string{
			"Set JIRA_INSTANCE_URL environment variable",
			"Set JIRA_USERNAME environment variable",
			"Set JIRA_API_TOKEN environment variable",
			"Run 'mcq config setup' for guided setup",
		},
	}

	ModelNotAvailableError = &UserError{
		Code:       "MODEL_NOT_AVAILABLE",
		Message:    "❌ AI Model Not Available",
		Suggestion: "Set up API keys for your preferred model",
		Troubleshoot: []string{
			"Set ANTHROPIC_API_KEY for Claude models",
			"Set OPENAI_API_KEY for GPT models",
			"Run 'mcq ai models' to see available options",
			"Run 'mcq config setup' for guided setup",
		},
	}

	IssueNotFoundError = &UserError{
		Code:       "ISSUE_NOT_FOUND",
		Message:    "❌ JIRA Issue Not Found",
		Suggestion: "Check the issue key and try again",
		Troubleshoot: []string{
			"Verify the issue key is correct (e.g., PROJ-123)",
			"Check if you have access to this issue",
			"Verify the JIRA project prefix is correct",
			"Run 'mcq jira list' to see available issues",
		},
	}

	ContextGatheringError = &UserError{
		Code:       "CONTEXT_GATHERING_FAILED",
		Message:    "⚠️  Context Gathering Failed",
		Suggestion: "Continuing without context (results may be less accurate)",
		Troubleshoot: []string{
			"Check if you're in a Git repository",
			"Verify file permissions for README and config files",
			"Use --no-context flag to skip context gathering",
			"Run 'mcq context test' to diagnose issues",
		},
	}

	ClipboardError = &UserError{
		Code:       "CLIPBOARD_FAILED",
		Message:    "⚠️  Clipboard Copy Failed",
		Suggestion: "Content is still displayed above",
		Troubleshoot: []string{
			"Check if pbcopy is available on macOS",
			"Try copying the content manually",
			"Use --no-clipboard flag to skip clipboard copy",
		},
	}
)

// NewUserError creates a new user error with original error
func NewUserError(userErr *UserError, originalErr error) *UserError {
	return &UserError{
		Code:         userErr.Code,
		Message:      userErr.Message,
		Suggestion:   userErr.Suggestion,
		Troubleshoot: userErr.Troubleshoot,
		OriginalErr:  originalErr,
	}
}

// WrapError wraps a technical error with user-friendly context
func WrapError(originalErr error, context string) *UserError {
	// Check if it's already a UserError
	if userErr, ok := originalErr.(*UserError); ok {
		return userErr
	}

	// Map common technical errors to user-friendly ones
	errStr := strings.ToLower(originalErr.Error())

	switch {
	case strings.Contains(errStr, "authentication") || strings.Contains(errStr, "unauthorized"):
		return NewUserError(JiraAuthError, originalErr)
	case strings.Contains(errStr, "not found") || strings.Contains(errStr, "404"):
		return NewUserError(IssueNotFoundError, originalErr)
	case strings.Contains(errStr, "api key") || strings.Contains(errStr, "token"):
		return NewUserError(ModelNotAvailableError, originalErr)
	case strings.Contains(errStr, "configuration") || strings.Contains(errStr, "config"):
		return NewUserError(JiraConfigError, originalErr)
	// Note: Context gathering errors are handled gracefully in the calling code
	// We don't auto-classify generic "context" or "repository" errors here
	case strings.Contains(errStr, "clipboard") || strings.Contains(errStr, "pbcopy"):
		return NewUserError(ClipboardError, originalErr)
	default:
		// Generic user error for unknown issues
		return &UserError{
			Code:       "UNKNOWN_ERROR",
			Message:    fmt.Sprintf("❌ %s", context),
			Suggestion: "Please try again or check your configuration",
			Troubleshoot: []string{
				"Run 'mcq config test' to validate setup",
				"Use --verbose flag for detailed error information",
				"Check the documentation for troubleshooting",
			},
			OriginalErr: originalErr,
		}
	}
}
