package ai

import "fmt"

// GenerateUserStory generates a user story using the specified model
func GenerateUserStory(model ModelConfig, featureRequest string, repoContext *RepoContext) (string, error) {
	showProgress(model, featureRequest)

	if model.Provider == "anthropic" {
		result, err := generateUserStoryClaude(model.APIKey, featureRequest, repoContext)
		if err != nil {
			fmt.Printf("\n⚠️  Claude API error: %v\n", err)
		}
		return result, err
	}
	result, err := generateUserStoryOpenAI(model.APIKey, featureRequest, model.ModelID, repoContext)
	if err != nil {
		fmt.Printf("\n⚠️  OpenAI API error: %v\n", err)
	}
	return result, err
}

// showProgress displays progress indicators
func showProgress(model ModelConfig, featureRequest string) {
	fmt.Printf("🤖 Generating user story with %s...\n", model.Name)
	fmt.Printf("📝 Feature request: %s\n\n", featureRequest)
}
