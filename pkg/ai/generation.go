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
