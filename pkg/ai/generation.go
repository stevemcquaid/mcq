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

// GenerateImprovedDescription generates an improved version of a description using the specified model
func GenerateImprovedDescription(model ModelConfig, originalDescription string, repoContext *RepoContext) (string, error) {
	showImprovementProgress(model)

	if model.Provider == "anthropic" {
		result, err := generateImprovedDescriptionClaude(model.APIKey, originalDescription, repoContext)
		if err != nil {
			fmt.Printf("\n⚠️  Claude API error: %v\n", err)
		}
		return result, err
	}
	result, err := generateImprovedDescriptionOpenAI(model.APIKey, originalDescription, model.ModelID, repoContext)
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

// showImprovementProgress displays progress for description improvement
func showImprovementProgress(model ModelConfig) {
	fmt.Printf("🤖 Improving description with %s...\n", model.Name)
}

// GenerateDescriptionFromTitle generates a description from a Jira issue title
func GenerateDescriptionFromTitle(model ModelConfig, title string, repoContext *RepoContext) (string, error) {
	fmt.Printf("🤖 Generating description from title with %s...\n", model.Name)

	if model.Provider == "anthropic" {
		result, err := generateDescriptionFromTitleClaude(model.APIKey, title, repoContext)
		if err != nil {
			fmt.Printf("\n⚠️  Claude API error: %v\n", err)
		}
		return result, err
	}
	result, err := generateDescriptionFromTitleOpenAI(model.APIKey, title, model.ModelID, repoContext)
	if err != nil {
		fmt.Printf("\n⚠️  OpenAI API error: %v\n", err)
	}
	return result, err
}
