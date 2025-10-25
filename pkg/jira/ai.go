package jira

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/ai"
)

// AIExtractor handles AI-powered title extraction for JIRA issues
type AIExtractor struct {
	// This will be set by the commands package
	SelectModelFunc func(string) (ai.ModelConfig, error)
}

// NewAIExtractor creates a new AI extractor
func NewAIExtractor(selectModelFunc func(string) (ai.ModelConfig, error)) *AIExtractor {
	return &AIExtractor{
		SelectModelFunc: selectModelFunc,
	}
}

// ExtractTitleWithAI uses AI to extract a concise title from the user story
func (ae *AIExtractor) ExtractTitleWithAI(userStory, featureRequest string) (string, error) {
	if ae.SelectModelFunc == nil {
		return "", fmt.Errorf("AI functionality not available - SelectModelFunc not set")
	}

	// Get the model configuration
	model, err := ae.SelectModelFunc("") // Use auto-detection
	if err != nil {
		return "", fmt.Errorf("failed to select model: %w", err)
	}

	// Create a focused prompt for title extraction
	config := ai.GetTitleExtractionPromptConfig(featureRequest, userStory)
	prompt := ai.GeneratePrompt(config)

	// Use the AI package to generate the title
	title, err := ai.GenerateUserStory(model, prompt, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate title: %w", err)
	}

	// Clean up the title
	title = cleanTitle(title)
	return title, nil
}

// cleanTitle cleans and truncates the title to fit JIRA requirements
func cleanTitle(title string) string {
	// Remove any extra whitespace
	title = strings.TrimSpace(title)

	// Remove any markdown formatting
	title = strings.TrimPrefix(title, "#")
	title = strings.TrimPrefix(title, "##")
	title = strings.TrimPrefix(title, "###")
	title = strings.TrimSpace(title)

	// Truncate if too long
	if len(title) > 100 {
		title = title[:97] + "..."
	}

	return title
}
