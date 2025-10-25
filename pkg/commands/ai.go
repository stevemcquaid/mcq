package commands

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// ModelConfig represents configuration for an AI model
type ModelConfig struct {
	Name        string
	Provider    string
	APIKey      string
	ModelID     string
	Description string
}

// ModelProvider represents different AI providers
type ModelProvider string

const (
	ProviderAnthropic ModelProvider = "anthropic"
	ProviderOpenAI    ModelProvider = "openai"
)

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Logger *slog.Logger
}

// Global logger configuration
var loggerConfig LoggerConfig

// Log levels mapping
var (
	LogLevelOff      = slog.LevelError - 1 // Custom level for no output
	LogLevelBasic    = slog.LevelInfo
	LogLevelDetailed = slog.LevelDebug
	LogLevelVerbose  = slog.LevelDebug - 1 // Custom level for verbose output
)

// Constants
const (
	DefaultMaxTokens = 4000
	GPT5Prefix       = "gpt-5"
)

// Available models configuration
var (
	models = map[string]ModelConfig{
		"claude": {
			Name:        "Claude Sonnet 4.5",
			Provider:    string(ProviderAnthropic),
			ModelID:     "claude-sonnet-4-5-20250929",
			Description: "Latest Claude model for complex reasoning",
		},
		"gpt-4o": {
			Name:        "GPT-4o",
			Provider:    string(ProviderOpenAI),
			ModelID:     "gpt-4o",
			Description: "Previous generation GPT model",
		},
		"gpt-5": {
			Name:        "GPT-5",
			Provider:    string(ProviderOpenAI),
			ModelID:     "gpt-5",
			Description: "Full power, best for complex tasks",
		},
		"gpt-5-mini": {
			Name:        "GPT-5 Mini",
			Provider:    string(ProviderOpenAI),
			ModelID:     "gpt-5-mini",
			Description: "Faster and more cost-effective",
		},
		"gpt-5-nano": {
			Name:        "GPT-5 Nano",
			Provider:    string(ProviderOpenAI),
			ModelID:     "gpt-5-nano",
			Description: "Optimized for simple tasks",
		},
	}

	// Model selection order for interactive choice
	modelOrder = []string{"claude", "gpt-4o", "gpt-5", "gpt-5-mini", "gpt-5-nano"}
)

// setupLogger configures the logger based on verbosity level
func setupLogger(verbosityLevel int) {
	var level slog.Level
	switch verbosityLevel {
	case 0:
		level = LogLevelOff
	case 1:
		level = LogLevelBasic
	case 2:
		level = LogLevelDetailed
	case 3:
		level = LogLevelVerbose
	default:
		level = LogLevelOff
	}

	// Create a custom handler that respects our verbosity levels
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize the output format
			if a.Key == slog.TimeKey {
				return slog.Attr{} // Remove timestamp
			}
			return a
		},
	})

	loggerConfig.Logger = slog.New(handler)
}

// logBasic logs basic process information
func logBasic(msg string, args ...interface{}) {
	if loggerConfig.Logger != nil {
		loggerConfig.Logger.Info(msg, args...)
	}
}

// logDetailed logs detailed information
func logDetailed(msg string, args ...interface{}) {
	if loggerConfig.Logger != nil {
		loggerConfig.Logger.Debug(msg, args...)
	}
}

// logVerbose logs verbose information (streaming details)
func logVerbose(msg string, args ...interface{}) {
	if loggerConfig.Logger != nil {
		// Use a custom level for verbose output
		loggerConfig.Logger.Log(context.Background(), LogLevelVerbose, msg, args...)
	}
}

// logAPI logs API-related information
func logAPI(operation, details string) {
	if loggerConfig.Logger != nil {
		loggerConfig.Logger.Debug("API operation", "operation", operation, "details", details)
	}
}

// logError logs error information
func logError(operation string, err error) {
	if loggerConfig.Logger != nil {
		loggerConfig.Logger.Error("Operation failed", "operation", operation, "error", err)
	}
}

// AIJira converts a vague feature request to a user story and copies it to clipboard
func AIJira(args []string, modelFlag string, verbosityLevel int) error {
	// Setup logger based on verbosity level
	setupLogger(verbosityLevel)

	featureRequest := strings.Join(args, " ")
	logBasic("Starting AIJira", "feature_request", featureRequest)
	logBasic("Configuration", "model_flag", modelFlag, "verbosity_level", verbosityLevel)

	// Get API keys
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")

	logDetailed("API Keys",
		"anthropic", func() string {
			if anthropicAPIKey != "" {
				return "***" + anthropicAPIKey[len(anthropicAPIKey)-4:]
			}
			return "not set"
		}(),
		"openai", func() string {
			if openaiAPIKey != "" {
				return "***" + openaiAPIKey[len(openaiAPIKey)-4:]
			}
			return "not set"
		}())

	// Select model
	logBasic("Selecting model")
	selectedModel, err := selectModel(modelFlag, anthropicAPIKey, openaiAPIKey)
	if err != nil {
		logError("model selection", err)
		return err
	}
	logBasic("Selected model", "name", selectedModel.Name, "provider", selectedModel.Provider)

	// Show progress
	showProgress(selectedModel, featureRequest)

	// Generate user story
	logBasic("Generating user story")
	userStory, err := generateUserStory(selectedModel, featureRequest)
	if err != nil {
		logError("user story generation", err)
		return fmt.Errorf("failed to generate user story: %w", err)
	}
	logBasic("User story generated successfully", "length", len(userStory))

	// Copy to clipboard and display result
	logBasic("Copying to clipboard and displaying result")
	return copyAndDisplayResult(userStory)
}

// selectModel determines which AI model to use
func selectModel(modelFlag, anthropicAPIKey, openaiAPIKey string) (ModelConfig, error) {
	logDetailed("Model selection", "flag", modelFlag)

	// If model is explicitly specified, validate and return it
	if modelFlag != "" {
		logDetailed("Explicit model specified", "model", modelFlag)
		model, exists := models[modelFlag]
		if !exists {
			err := fmt.Errorf("unsupported model: %s", modelFlag)
			logError("model validation", err)
			return ModelConfig{}, err
		}
		logDetailed("Model found in configuration", "name", model.Name)

		// Set API key based on provider
		if model.Provider == string(ProviderAnthropic) {
			logDetailed("Model requires Anthropic API key")
			if anthropicAPIKey == "" {
				err := fmt.Errorf("ANTHROPIC_API_KEY is required for Claude model")
				logError("API key validation", err)
				return ModelConfig{}, err
			}
			model.APIKey = anthropicAPIKey
			logDetailed("Anthropic API key set")
		} else if model.Provider == string(ProviderOpenAI) {
			logDetailed("Model requires OpenAI API key")
			if openaiAPIKey == "" {
				err := fmt.Errorf("OPENAI_API_KEY is required for %s model", model.Name)
				logError("API key validation", err)
				return ModelConfig{}, err
			}
			model.APIKey = openaiAPIKey
			logDetailed("OpenAI API key set")
		}

		logDetailed("Model selection complete", "name", model.Name)
		return model, nil
	}

	// Auto-detect based on available API keys
	hasAnthropic := anthropicAPIKey != ""
	hasOpenAI := openaiAPIKey != ""

	logDetailed("Auto-detection", "anthropic_available", hasAnthropic, "openai_available", hasOpenAI)

	if !hasAnthropic && !hasOpenAI {
		err := fmt.Errorf("no API keys found. Please set either ANTHROPIC_API_KEY or OPENAI_API_KEY")
		logError("API key detection", err)
		return ModelConfig{}, err
	}

	if hasAnthropic && hasOpenAI {
		logDetailed("Both API keys available, prompting user for selection")
		return interactiveModelSelection(anthropicAPIKey, openaiAPIKey)
	}

	// Only one provider available
	if hasAnthropic {
		logDetailed("Only Anthropic API key available, selecting Claude")
		model := models["claude"]
		model.APIKey = anthropicAPIKey
		return model, nil
	}

	// Default to GPT-5 for OpenAI
	logDetailed("Only OpenAI API key available, selecting GPT-5")
	model := models["gpt-5"]
	model.APIKey = openaiAPIKey
	return model, nil
}

// interactiveModelSelection handles user choice when both API keys are available
func interactiveModelSelection(anthropicAPIKey, openaiAPIKey string) (ModelConfig, error) {
	logDetailed("Starting interactive model selection")
	fmt.Println("🔑 Both Claude and OpenAI API keys are available.")
	fmt.Println("Which model would you like to use?")

	availableModels := 0
	for i, modelKey := range modelOrder {
		model := models[modelKey]
		// Check if API key is available for this model
		if (model.Provider == string(ProviderAnthropic) && anthropicAPIKey != "") ||
			(model.Provider == string(ProviderOpenAI) && openaiAPIKey != "") {
			fmt.Printf("%d. %s (%s) - %s\n", i+1, model.Name, strings.ToUpper(model.Provider[:1])+model.Provider[1:], model.Description)
			availableModels++
			logDetailed("Available model", "index", i+1, "name", model.Name, "provider", model.Provider)
		}
	}

	logDetailed("Total available models", "count", availableModels)
	fmt.Print("Enter choice (1-5): ")

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		logError("user input", err)
		return ModelConfig{}, fmt.Errorf("invalid input: %w", err)
	}

	logDetailed("User selected choice", "choice", choice)

	if choice < 1 || choice > len(modelOrder) {
		err := fmt.Errorf("choice %d out of range", choice)
		logError("choice validation", err)
		return ModelConfig{}, fmt.Errorf("invalid choice. Please select 1-5")
	}

	selectedModelKey := modelOrder[choice-1]
	model := models[selectedModelKey]
	logDetailed("Selected model key", "key", selectedModelKey)

	// Set appropriate API key
	if model.Provider == string(ProviderAnthropic) {
		model.APIKey = anthropicAPIKey
		logDetailed("Set Anthropic API key for model")
	} else {
		model.APIKey = openaiAPIKey
		logDetailed("Set OpenAI API key for model")
	}

	logDetailed("Interactive selection complete", "name", model.Name)
	return model, nil
}

// showProgress displays progress indicators
func showProgress(model ModelConfig, featureRequest string) {
	fmt.Printf("🤖 Generating user story with %s...\n", model.Name)
	fmt.Printf("📝 Feature request: %s\n\n", featureRequest)
}

// showConnectionProgressWithModel displays progress during API connection setup with model info
func showConnectionProgressWithModel(provider, modelName string) {
	fmt.Printf("🔌 Connecting to %s API (%s)...\n", provider, modelName)
}

// showStreamingProgress displays progress when streaming starts
func showStreamingProgress() {
	fmt.Print("💭 ")
}

// getModelDisplayName returns a user-friendly name for the model ID
func getModelDisplayName(modelID string) string {
	switch modelID {
	case "gpt-4o":
		return "GPT-4o"
	case "gpt-5":
		return "GPT-5"
	case "gpt-5-mini":
		return "GPT-5 Mini"
	case "gpt-5-nano":
		return "GPT-5 Nano"
	default:
		return modelID
	}
}

// generateUserStory generates a user story using the specified model
func generateUserStory(model ModelConfig, featureRequest string) (string, error) {
	if model.Provider == string(ProviderAnthropic) {
		return generateUserStoryClaude(model.APIKey, featureRequest)
	}
	return generateUserStoryOpenAI(model.APIKey, featureRequest, model.ModelID)
}

// copyAndDisplayResult copies the user story to clipboard and displays it
func copyAndDisplayResult(userStory string) error {
	fmt.Println("\n📋 Copying to clipboard...")

	if err := copyToClipboard(userStory); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	fmt.Println("✅ User story generated and copied to clipboard!")
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println(userStory)
	fmt.Println(strings.Repeat("=", 60))
	return nil
}

// generateUserStoryClaude calls the Anthropic API with streaming
func generateUserStoryClaude(apiKey, featureRequest string) (string, error) {
	logBasic("Starting Claude API request")
	showConnectionProgressWithModel("Anthropic", "Claude Sonnet 4.5")
	prompt := createPrompt(featureRequest)
	logDetailed("Prompt created", "length", len(prompt))

	request := AnthropicRequest{
		Model:     "claude-sonnet-4-5-20250929",
		MaxTokens: 4000,
		Stream:    true,
		Messages: []AnthropicMessage{
			{Role: "user", Content: prompt},
		},
	}
	logAPI("Claude Request", fmt.Sprintf("Model: %s, MaxTokens: %d, Stream: %t", request.Model, request.MaxTokens, request.Stream))

	jsonData, err := json.Marshal(request)
	if err != nil {
		logError("JSON marshaling", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	logDetailed("Request JSON created", "size", len(jsonData))

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		logError("HTTP request creation", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	logAPI("Claude Headers", "Content-Type: application/json, anthropic-version: 2023-06-01")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	logDetailed("Making HTTP request to Claude API")
	resp, err := client.Do(req)
	if err != nil {
		logError("HTTP request execution", err)
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			if loggerConfig.Logger != nil {
				loggerConfig.Logger.Error("Operation failed", "operation", "close response body", "error", err)
			}
		}
	}()

	logAPI("Claude Response", fmt.Sprintf("Status: %d %s", resp.StatusCode, resp.Status))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
		logError("Claude API Error", err)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	logBasic("Claude API request successful, processing stream")
	showStreamingProgress()
	return processClaudeStream(resp.Body)
}

// generateUserStoryOpenAI calls the OpenAI API with streaming
func generateUserStoryOpenAI(apiKey, featureRequest, modelID string) (string, error) {
	logBasic("Starting OpenAI API request", "model", modelID)
	// Get the model name for display
	modelName := getModelDisplayName(modelID)
	showConnectionProgressWithModel("OpenAI", modelName)
	prompt := createPrompt(featureRequest)
	logDetailed("Prompt created", "length", len(prompt))

	client := openai.NewClient(apiKey)
	req := createOpenAIRequest(modelID, prompt)

	logDetailed("Creating OpenAI streaming request")
	// Create a context with timeout for stream creation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		logError("OpenAI stream creation", err)
		return "", fmt.Errorf("failed to create streaming request: %w", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			if loggerConfig.Logger != nil {
				loggerConfig.Logger.Error("Operation failed", "operation", "close stream", "error", err)
			}
		}
	}()

	logBasic("OpenAI stream created successfully, processing")
	showStreamingProgress()
	return processOpenAIStream(stream)
}

// createOpenAIRequest creates the appropriate OpenAI request based on model type
func createOpenAIRequest(modelID, prompt string) openai.ChatCompletionRequest {
	baseRequest := openai.ChatCompletionRequest{
		Model: modelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Stream: true,
	}

	if strings.HasPrefix(modelID, GPT5Prefix) {
		// For GPT-5 models, omit max tokens parameter
		logAPI("OpenAI Request", fmt.Sprintf("Model: %s (GPT-5 - no max tokens), Stream: %t", baseRequest.Model, baseRequest.Stream))
		return baseRequest
	}

	// For older models, add max tokens
	baseRequest.MaxTokens = DefaultMaxTokens
	logAPI("OpenAI Request", fmt.Sprintf("Model: %s, MaxTokens: %d, Stream: %t", baseRequest.Model, baseRequest.MaxTokens, baseRequest.Stream))
	return baseRequest
}

// createPrompt creates the standardized prompt for user story generation
func createPrompt(featureRequest string) string {
	return fmt.Sprintf(`Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations.

Feature Request: %s

Please provide a comprehensive user story with:
1. The main user story in the specified format
2. Acceptance criteria
3. Any relevant technical notes or considerations`, featureRequest)
}

// processClaudeStream processes the streaming response from Claude
func processClaudeStream(body io.ReadCloser) (string, error) {
	logBasic("Processing Claude stream")
	var fullResponse strings.Builder
	scanner := bufio.NewScanner(body)

	eventCount := 0
	contentBlocks := 0

	for scanner.Scan() {
		line := scanner.Text()
		logVerbose("Stream line", "line", line)

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			logVerbose("Data event", "data", data)

			if data == "[DONE]" {
				logBasic("Stream completed with [DONE] marker")
				break
			}

			var event StreamingEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				logError("JSON unmarshal", err)
				continue
			}

			eventCount++
			logVerbose("Event", "count", eventCount, "type", event.Type)

			if event.Type == "content_block_delta" && event.Delta.Text != "" {
				contentBlocks++
				logVerbose("Content delta", "count", contentBlocks, "text", event.Delta.Text)
				fmt.Print(event.Delta.Text)
				fullResponse.WriteString(event.Delta.Text)
			}
		}
	}

	fmt.Println()
	logDetailed("Stream processing complete", "events", eventCount, "content_blocks", contentBlocks)

	if err := scanner.Err(); err != nil {
		logError("stream scanning", err)
		return "", fmt.Errorf("error reading stream: %w", err)
	}

	response := fullResponse.String()
	logDetailed("Final response", "length", len(response))

	if response == "" {
		err := fmt.Errorf("no content in response")
		logError("response validation", err)
		return "", err
	}

	return response, nil
}

// processOpenAIStream processes the streaming response from OpenAI
func processOpenAIStream(stream *openai.ChatCompletionStream) (string, error) {
	logBasic("Processing OpenAI stream")
	var fullResponse strings.Builder

	chunkCount := 0
	contentChunks := 0

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logBasic("Stream completed with EOF")
				break
			}
			logError("stream receive", err)
			return "", fmt.Errorf("stream error: %w", err)
		}

		chunkCount++
		logVerbose("Received chunk", "count", chunkCount)

		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			content := response.Choices[0].Delta.Content
			contentChunks++
			logVerbose("Content chunk", "count", contentChunks, "text", content)
			fmt.Print(content)
			fullResponse.WriteString(content)
		} else {
			logVerbose("Empty chunk", "count", chunkCount)
		}
	}

	fmt.Println()
	logDetailed("Stream processing complete", "chunks", chunkCount, "content_chunks", contentChunks)

	response := fullResponse.String()
	logDetailed("Final response", "length", len(response))

	if response == "" {
		err := fmt.Errorf("no content in response")
		logError("response validation", err)
		return "", err
	}

	return response, nil
}

// copyToClipboard copies text to clipboard on macOS
func copyToClipboard(text string) error {
	logDetailed("Copying to clipboard", "length", len(text))
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)

	err := cmd.Run()
	if err != nil {
		logError("clipboard copy", err)
		return err
	}

	logBasic("Successfully copied to clipboard")
	return nil
}

// AnthropicMessage represents a message in the conversation
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicRequest represents the request to Anthropic API
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
	Stream    bool               `json:"stream"`
}

// StreamingEvent represents a streaming event from Anthropic API
type StreamingEvent struct {
	Type  string `json:"type"`
	Delta struct {
		Text string `json:"text"`
	} `json:"delta"`
	ContentBlock struct {
		Text string `json:"text"`
	} `json:"content_block"`
}
