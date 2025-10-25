// Package commands provides Jira integration functionality for the mcq CLI tool.
// This file contains the main JIRA command handlers and orchestrates the refactored components.
package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/ai"
	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/jira"
	"github.com/stevemcquaid/mcq/pkg/logger"
)

// askForConfirmation prompts the user for confirmation with a default value.
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

// copyToClipboard copies text to clipboard on macOS
func copyToClipboard(text string) error {
	logger.LogDetailed("Copying to clipboard", "length", len(text))
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)

	err := cmd.Run()
	if err != nil {
		logger.LogError("clipboard copy", err)
		return err
	}

	logger.LogBasic("Successfully copied to clipboard")
	return nil
}

// ShowJiraIssue displays detailed information about a Jira issue.
// This is the main entry point for the "mcq jira show" command.
func ShowJiraIssue(issueKey string) {
	manager, err := jira.NewManager()
	if err != nil {
		userErr := errors.WrapError(err, "Failed to create Jira manager")
		userErr.Display()
		return
	}

	if err := manager.ShowIssue(issueKey); err != nil {
		userErr := errors.WrapError(err, "Failed to show issue")
		userErr.Display()
		return
	}
}

// JiraNew creates a new Jira issue from a vague user story using AI
func JiraNew(args []string, modelFlag string, verbosityLevel int, contextConfig ai.ContextConfig) error {
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

	// Create Jira manager
	manager, err := jira.NewManager()
	if err != nil {
		userErr := errors.WrapError(err, "Failed to create Jira manager")
		userErr.Display()
		return userErr
	}

	// Set up AI extractor
	aiExtractor := jira.NewAIExtractor(ai.SelectModel)
	manager.SetAIExtractor(aiExtractor)

	// Create the Jira issue
	issueKey, err := manager.CreateIssue(userStory, featureRequest)
	if err != nil {
		userErr := errors.WrapError(err, "Failed to create Jira issue")
		userErr.Display()
		return userErr
	}

	// Display success message
	fmt.Printf("\n✅ Jira issue created successfully: %s\n", issueKey)
	fmt.Printf("🔗 You can view it at: %s/browse/%s\n", manager.GetBaseURL(), issueKey)

	return nil
}

// generateUserStoryForJira generates a user story using AI and returns it without copying to clipboard
func generateUserStoryForJira(featureRequest string, modelFlag string, verbosityLevel int, contextConfig ai.ContextConfig) (string, error) {
	// Set up logger
	logger.SetupLogger(verbosityLevel)
	logger.LogBasic("Starting JiraNew", "feature_request", featureRequest)

	// Gather repository context
	repoContext := ai.GatherContextIfNeeded(contextConfig)

	// Select and configure model
	selectedModel, err := ai.SelectModel(modelFlag)
	if err != nil {
		return "", err
	}

	// Generate user story
	userStory, err := ai.GenerateUserStory(selectedModel, featureRequest, repoContext)
	if err != nil {
		return "", fmt.Errorf("failed to generate user story: %w", err)
	}

	// Copy to clipboard (as requested)
	if err := copyToClipboard(userStory); err != nil {
		logger.LogError("clipboard copy", err)
		// Don't fail the entire operation if clipboard copy fails
		userErr := errors.WrapError(err, "Clipboard copy failed")
		userErr.Display()
	} else {
		fmt.Println("📋 User story copied to clipboard!")
	}

	return userStory, nil
}
