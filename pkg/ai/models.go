package ai

// Available models
var models = map[string]ModelConfig{
	"claude": {
		Name:        "Claude Sonnet 4.5",
		Provider:    "anthropic",
		ModelID:     "claude-sonnet-4-5-20250929",
		Description: "Latest Claude model for complex reasoning",
	},
	"gpt-4o": {
		Name:        "GPT-4o",
		Provider:    "openai",
		ModelID:     "gpt-4o",
		Description: "Previous generation GPT model",
	},
	"gpt-5": {
		Name:        "GPT-5",
		Provider:    "openai",
		ModelID:     "gpt-5",
		Description: "Full power, best for complex tasks",
	},
	"gpt-5-mini": {
		Name:        "GPT-5 Mini",
		Provider:    "openai",
		ModelID:     "gpt-5-mini",
		Description: "Faster and more cost-effective",
	},
	"gpt-5-nano": {
		Name:        "GPT-5 Nano",
		Provider:    "openai",
		ModelID:     "gpt-5-nano",
		Description: "Optimized for simple tasks",
	},
}

var modelOrder = []string{"claude", "gpt-4o", "gpt-5", "gpt-5-mini", "gpt-5-nano"}

// GetAvailableModels returns a list of available model names
func GetAvailableModels() []string {
	return modelOrder
}

// GetModel returns a model configuration by name
func GetModel(name string) (ModelConfig, bool) {
	model, exists := models[name]
	return model, exists
}
