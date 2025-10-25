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
