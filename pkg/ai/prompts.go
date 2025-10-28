package ai

import (
	"fmt"

	"github.com/stevemcquaid/mcq/pkg/logger"
)

// PromptType represents different types of AI prompts
type PromptType string

const (
	// PromptTypeUserStory generates user stories from feature requests
	PromptTypeUserStory PromptType = "user_story"
	// PromptTypeTitleExtraction extracts concise titles from user stories
	PromptTypeTitleExtraction PromptType = "title_extraction"
)

// PromptConfig holds configuration for prompt generation
type PromptConfig struct {
	Type              PromptType
	FeatureRequest    string
	UserStory         string
	RepositoryContext *RepoContext
}

// GeneratePrompt creates a standardized prompt based on the type and configuration
func GeneratePrompt(config PromptConfig) string {
	tm := GetTemplateManager()

	// Load templates if not already loaded
	if len(tm.templates) == 0 {
		if err := tm.LoadTemplates(); err != nil {
			logger.LogError("Failed to load templates", err)
			return getDefaultPrompt(config)
		}
	}

	// Prepare template data
	data := TemplateData{
		FeatureRequest:    config.FeatureRequest,
		UserStory:         config.UserStory,
		RepositoryContext: config.RepositoryContext,
	}

	// Generate prompt using template system
	prompt, err := tm.GeneratePromptFromTemplate(config.Type, data)
	if err != nil {
		logger.LogError("Failed to generate prompt from template", err)
		return getDefaultPrompt(config)
	}

	// Log prompt size for debugging
	logger.LogBasic("Generated prompt", "size_chars", len(prompt))
	if len(prompt) > 100000 {
		logger.LogBasic("Warning: Large prompt may exceed token limits", "size_chars", len(prompt))
	}

	return prompt
}

// getDefaultPrompt provides default prompts when templates fail
func getDefaultPrompt(config PromptConfig) string {
	switch config.Type {
	case PromptTypeUserStory:
		return createUserStoryPrompt(config.FeatureRequest, config.RepositoryContext)
	case PromptTypeTitleExtraction:
		return createTitleExtractionPrompt(config.FeatureRequest, config.UserStory)
	default:
		return ""
	}
}

// createUserStoryPrompt creates the standardized prompt for user story generation
func createUserStoryPrompt(featureRequest string, repoContext *RepoContext) string {
	basePrompt := `Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations. Provide ONLY the user story. 


Please provide a comprehensive user story:
1. With the main user story in the specified format
2. With acceptance criteria
3. With any relevant technical notes or considerations
4. Keep the total output under 1000 words

Do NOT add any additional questions or commentary. 
The response must ONLY be the user story. 
NOTHING ELSE.

Feature Request: %s
`

	// Add repository context if available
	if repoContext != nil {
		contextInfo := formatContextForPrompt(repoContext)
		basePrompt += contextInfo
	}

	return fmt.Sprintf(basePrompt, featureRequest)
}

// createTitleExtractionPrompt creates a focused prompt for title extraction
func createTitleExtractionPrompt(featureRequest, userStory string) string {
	return fmt.Sprintf(`Create a NEW concise, clear title (maximum 100 characters) for a Jira issue from the following user story and old title. The new title should be action-oriented and summarize the main goal or feature.
Provide ONLY the new jira title
Do NOT provide any other output.

Original Feature Request: %s

User Story: 
%s`, featureRequest, userStory)
}

// GetUserStoryPromptConfig creates a prompt configuration for user story generation
func GetUserStoryPromptConfig(featureRequest string, repoContext *RepoContext) PromptConfig {
	return PromptConfig{
		Type:              PromptTypeUserStory,
		FeatureRequest:    featureRequest,
		RepositoryContext: repoContext,
	}
}

// GetTitleExtractionPromptConfig creates a prompt configuration for title extraction
func GetTitleExtractionPromptConfig(featureRequest, userStory string) PromptConfig {
	return PromptConfig{
		Type:           PromptTypeTitleExtraction,
		FeatureRequest: featureRequest,
		UserStory:      userStory,
	}
}
