package ai

import "fmt"

// GenerateUserStory generates a user story using the specified model
func GenerateUserStory(model ModelConfig, featureRequest string, repoContext *RepoContext) (string, error) {
	showProgress(model, featureRequest)

	if model.Provider == "anthropic" {
		return generateUserStoryClaude(model.APIKey, featureRequest, repoContext)
	}
	return generateUserStoryOpenAI(model.APIKey, featureRequest, model.ModelID, repoContext)
}

// showProgress displays progress indicators
func showProgress(model ModelConfig, featureRequest string) {
	fmt.Printf("🤖 Generating user story with %s...\n", model.Name)
	fmt.Printf("📝 Feature request: %s\n\n", featureRequest)
}

// createPrompt creates the standardized prompt for user story generation
func createPrompt(featureRequest string, repoContext *RepoContext) string {
	basePrompt := `Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations.

Feature Request: %s

Please provide a comprehensive user story with:
1. The main user story in the specified format
2. Acceptance criteria
3. Any relevant technical notes or considerations`

	// Add repository context if available
	if repoContext != nil {
		contextInfo := formatContextForPrompt(repoContext)
		basePrompt += contextInfo
	}

	return fmt.Sprintf(basePrompt, featureRequest)
}
