package ai

import (
	"fmt"
	"os"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/logger"
)

// SelectModel determines which AI model to use
func SelectModel(modelFlag string) (ModelConfig, error) {
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")

	logger.LogDetailed("API Keys",
		"anthropic", maskAPIKey(anthropicAPIKey),
		"openai", maskAPIKey(openaiAPIKey))

	// If model is explicitly specified, validate and return it
	if modelFlag != "" {
		model, err := selectExplicitModel(modelFlag, anthropicAPIKey, openaiAPIKey)
		if err != nil {
			return ModelConfig{}, errors.WrapError(err, "Failed to select AI model")
		}
		return model, nil
	}

	// Auto-detect based on available API keys
	model, err := selectModelByAvailability(anthropicAPIKey, openaiAPIKey)
	if err != nil {
		return ModelConfig{}, errors.WrapError(err, "Failed to auto-detect AI model")
	}
	return model, nil
}

// maskAPIKey returns a masked version of the API key for logging
func maskAPIKey(apiKey string) string {
	if apiKey == "" {
		return "not set"
	}
	return "***" + apiKey[len(apiKey)-4:]
}

// selectExplicitModel selects a model when explicitly specified
func selectExplicitModel(modelFlag, anthropicAPIKey, openaiAPIKey string) (ModelConfig, error) {
	model, exists := models[modelFlag]
	if !exists {
		return ModelConfig{}, fmt.Errorf("unsupported model: %s", modelFlag)
	}

	// Set API key based on provider
	switch model.Provider {
	case "anthropic":
		if anthropicAPIKey == "" {
			return ModelConfig{}, errors.ModelNotAvailableError
		}
		model.APIKey = anthropicAPIKey
	case "openai":
		if openaiAPIKey == "" {
			return ModelConfig{}, errors.ModelNotAvailableError
		}
		model.APIKey = openaiAPIKey
	}

	logger.LogBasic("Selected model", "name", model.Name, "provider", model.Provider)
	return model, nil
}

// selectModelByAvailability selects a model based on available API keys
func selectModelByAvailability(anthropicAPIKey, openaiAPIKey string) (ModelConfig, error) {
	hasAnthropic := anthropicAPIKey != ""
	hasOpenAI := openaiAPIKey != ""

	if !hasAnthropic && !hasOpenAI {
		return ModelConfig{}, errors.ModelNotAvailableError
	}

	if hasAnthropic && hasOpenAI {
		return interactiveModelSelection(anthropicAPIKey, openaiAPIKey)
	}

	// Only one provider available
	if hasAnthropic {
		model := models["claude"]
		model.APIKey = anthropicAPIKey
		return model, nil
	}

	// Default to GPT-5 for OpenAI
	model := models["gpt-5"]
	model.APIKey = openaiAPIKey
	return model, nil
}

// interactiveModelSelection handles user choice when both API keys are available
func interactiveModelSelection(anthropicAPIKey, openaiAPIKey string) (ModelConfig, error) {
	fmt.Println("🔑 Both Claude and OpenAI API keys are available.")
	fmt.Println("Which model would you like to use?")

	availableModels := 0
	for i, modelKey := range modelOrder {
		model := models[modelKey]
		if isModelAvailable(model, anthropicAPIKey, openaiAPIKey) {
			fmt.Printf("%d. %s (%s) - %s\n", i+1, model.Name, capitalize(model.Provider), model.Description)
			availableModels++
		}
	}

	fmt.Print("Enter choice (1-5): ")
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		fmt.Println("\n⚠️  Error reading input, using default model.")
		fmt.Println("   This is normal in non-interactive environments.")
		// Return a default model instead of error
		if anthropicAPIKey != "" {
			model := models["claude"]
			model.APIKey = anthropicAPIKey
			return model, nil
		}
		if openaiAPIKey != "" {
			model := models["gpt-5"]
			model.APIKey = openaiAPIKey
			return model, nil
		}
		return ModelConfig{}, errors.ModelNotAvailableError
	}

	if choice < 1 || choice > len(modelOrder) {
		return ModelConfig{}, fmt.Errorf("invalid choice. Please select 1-5")
	}

	selectedModelKey := modelOrder[choice-1]
	model := models[selectedModelKey]

	// Set appropriate API key
	if model.Provider == "anthropic" {
		model.APIKey = anthropicAPIKey
	} else {
		model.APIKey = openaiAPIKey
	}

	return model, nil
}

// isModelAvailable checks if a model is available with the given API keys
func isModelAvailable(model ModelConfig, anthropicAPIKey, openaiAPIKey string) bool {
	return (model.Provider == "anthropic" && anthropicAPIKey != "") ||
		(model.Provider == "openai" && openaiAPIKey != "")
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
