package ai

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/logger"
)

// handleOpenAIStreamError handles errors from CreateChatCompletionStream
// and provides helpful diagnostics
func handleOpenAIStreamError(err error, apiKey, modelID string) error {
	// Log the underlying error for debugging
	logger.LogError("OpenAI stream creation failed", err)

	errStr := strings.ToLower(err.Error())
	errMsg := err.Error()

	// Check for empty API key
	if apiKey == "" {
		return errors.WrapError(
			fmt.Errorf("OPENAI_API_KEY environment variable is not set"),
			"Failed to create streaming request",
		)
	}

	// Check for invalid model errors
	if strings.Contains(errMsg, "is not supported with this method") {
		return errors.WrapError(
			fmt.Errorf("model '%s' is not supported for chat completions: %v", modelID, err),
			"Failed to create streaming request",
		)
	}

	// Check for authentication errors
	if strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "invalid api key") ||
		strings.Contains(errStr, "incorrect api key") || strings.Contains(errStr, "authentication") {
		return errors.WrapError(
			fmt.Errorf("invalid OpenAI API key: %v", err),
			"Failed to create streaming request",
		)
	}

	// Check for token/context limit errors
	if strings.Contains(errStr, "context length") || strings.Contains(errStr, "token") ||
		strings.Contains(errStr, "maximum context") {
		logger.LogError("Token/context limit error detected", err)
		fmt.Printf("\n⚠️  Error: Context may be too large for the model\n")
		fmt.Printf("💡 Try reducing context with --no-context or specific context flags\n")
	}

	// Check for timeout errors
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return errors.WrapError(
			fmt.Errorf("request timed out: %v", err),
			"Failed to create streaming request",
		)
	}

	// Check for network errors
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "network") ||
		strings.Contains(errStr, "dial") {
		return errors.WrapError(
			fmt.Errorf("network error: %v", err),
			"Failed to create streaming request",
		)
	}

	// Return wrapped error with original error message preserved
	return errors.WrapError(err, "Failed to create streaming request")
}

// generateDescriptionFromTitleOpenAI generates a description from a title using OpenAI
func generateDescriptionFromTitleOpenAI(apiKey, title, modelID string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting OpenAI API request for description generation from title", "model", modelID)
	modelName := getModelDisplayName(modelID)
	showConnectionProgress("OpenAI", modelName)

	config := GetDescriptionFromTitlePromptConfig(title, repoContext)
	prompt := GeneratePrompt(config)
	client := openai.NewClient(apiKey)
	req := createOpenAIRequest(modelID, prompt)

	// Create a context with timeout for stream creation and processing
	// Use a longer timeout to allow for large responses
	ctx, cancel := context.WithTimeout(context.Background(), OpenAIStreamTimeout*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		cancel()
		return "", handleOpenAIStreamError(err, apiKey, modelID)
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			logger.LogError("close stream", closeErr)
		}
	}()

	logger.LogBasic("OpenAI stream created successfully, processing")
	showStreamingProgress()
	return processOpenAIStream(stream, ctx)
}

// generateImprovedDescriptionOpenAI generates an improved description using OpenAI
func generateImprovedDescriptionOpenAI(apiKey, originalDescription, modelID string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting OpenAI API request for description improvement", "model", modelID)
	modelName := getModelDisplayName(modelID)
	showConnectionProgress("OpenAI", modelName)

	config := GetDescriptionImprovementPromptConfig(originalDescription, repoContext)
	prompt := GeneratePrompt(config)
	client := openai.NewClient(apiKey)
	req := createOpenAIRequest(modelID, prompt)

	// Create a context with timeout for stream creation and processing
	// Use a longer timeout to allow for large responses
	ctx, cancel := context.WithTimeout(context.Background(), OpenAIStreamTimeout*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		cancel()
		return "", handleOpenAIStreamError(err, apiKey, modelID)
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			logger.LogError("close stream", closeErr)
		}
	}()

	logger.LogBasic("OpenAI stream created successfully, processing")
	showStreamingProgress()
	return processOpenAIStream(stream, ctx)
}

// generateUserStoryOpenAI calls the OpenAI API with streaming
func generateUserStoryOpenAI(apiKey, featureRequest, modelID string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting OpenAI API request", "model", modelID)
	modelName := getModelDisplayName(modelID)
	showConnectionProgress("OpenAI", modelName)

	config := GetUserStoryPromptConfig(featureRequest, repoContext)
	prompt := GeneratePrompt(config)
	client := openai.NewClient(apiKey)
	req := createOpenAIRequest(modelID, prompt)

	// Create a context with timeout for stream creation and processing
	// Use a longer timeout to allow for large responses
	ctx, cancel := context.WithTimeout(context.Background(), OpenAIStreamTimeout*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		cancel()
		return "", handleOpenAIStreamError(err, apiKey, modelID)
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			logger.LogError("close stream", closeErr)
		}
	}()

	logger.LogBasic("OpenAI stream created successfully, processing")
	showStreamingProgress()
	return processOpenAIStream(stream, ctx)
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
		return baseRequest
	}

	// For older models, add max tokens
	baseRequest.MaxTokens = DefaultMaxTokens
	return baseRequest
}

// handleOpenAIStreamRecvError handles errors that occur during stream processing
func handleOpenAIStreamRecvError(err error, partialResponse string) error {
	// Log the underlying error for debugging
	logger.LogError("OpenAI stream processing failed", err)

	errStr := strings.ToLower(err.Error())
	errMsg := err.Error()

	// Check for API error responses (these come from the stream itself)
	if strings.Contains(errMsg, "error, status code:") {
		// Extract status code if possible
		var statusCode string
		if strings.Contains(errMsg, "status code:") {
			parts := strings.Split(errMsg, "status code:")
			if len(parts) > 1 {
				statusCode = strings.Split(parts[1], ",")[0]
				statusCode = strings.TrimSpace(statusCode)
			}
		}

		// Extract message
		var message string
		if strings.Contains(errMsg, "message:") {
			parts := strings.Split(errMsg, "message:")
			if len(parts) > 1 {
				message = strings.TrimSpace(parts[1])
			}
		}

		if message != "" {
			// Check for specific error types
			if strings.Contains(errStr, "rate_limit") || strings.Contains(errStr, "rate limit") {
				return errors.WrapError(
					fmt.Errorf("rate limit exceeded: %s", message),
					"Stream error",
				)
			}
			if strings.Contains(errStr, "context_length") || strings.Contains(errStr, "token") {
				return errors.WrapError(
					fmt.Errorf("context/token limit exceeded: %s", message),
					"Stream error",
				)
			}
			if strings.Contains(errStr, "invalid_api_key") || strings.Contains(errStr, "authentication") {
				return errors.WrapError(
					fmt.Errorf("authentication error: %s", message),
					"Stream error",
				)
			}

			// Return error with status code and message
			if statusCode != "" {
				return errors.WrapError(
					fmt.Errorf("API error (status %s): %s", statusCode, message),
					"Stream error",
				)
			}
			return errors.WrapError(
				fmt.Errorf("API error: %s", message),
				"Stream error",
			)
		}
	}

	// Check for too many empty messages
	if strings.Contains(errMsg, "too many empty messages") {
		return errors.WrapError(
			fmt.Errorf("stream sent too many empty messages - possible connection issue: %v", err),
			"Stream error",
		)
	}

	// Check for network/connection errors
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "network") ||
		strings.Contains(errStr, "dial") || strings.Contains(errStr, "broken pipe") {
		if partialResponse != "" {
			logger.LogBasic("Partial response received before network error", "length", len(partialResponse))
		}
		return errors.WrapError(
			fmt.Errorf("network error during streaming: %v", err),
			"Stream error",
		)
	}

	// Check for timeout errors
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") ||
		strings.Contains(errStr, "context deadline exceeded") {
		if partialResponse != "" {
			logger.LogBasic("Partial response received before timeout", "length", len(partialResponse))
			fmt.Printf("\n⚠️  Request timed out after receiving partial response (%d characters)\n", len(partialResponse))
			fmt.Printf("💡 The response may be incomplete. You can try:\n")
			fmt.Printf("   • Reducing context size with --no-context flag\n")
			fmt.Printf("   • Using a faster model (e.g., gpt-5-mini instead of gpt-5)\n")
			fmt.Printf("   • Breaking your request into smaller parts\n")
		} else {
			fmt.Printf("\n⚠️  Request timed out (%d seconds)\n", OpenAIStreamTimeout)
			fmt.Printf("💡 The response may be taking longer than expected. Try:\n")
			fmt.Printf("   • Reducing context size with --no-context flag\n")
			fmt.Printf("   • Using a faster model\n")
		}
		return errors.WrapError(
			fmt.Errorf("request timed out during streaming: %v", err),
			"Stream error",
		)
	}

	// Check for unmarshaling errors
	if strings.Contains(errStr, "unmarshal") || strings.Contains(errStr, "json") {
		return errors.WrapError(
			fmt.Errorf("invalid response format from API: %v", err),
			"Stream error",
		)
	}

	// If we have partial content, log it but still return error
	if partialResponse != "" {
		logger.LogBasic("Partial response received before error", "length", len(partialResponse))
	}

	// Return wrapped error with original error message preserved
	return errors.WrapError(err, "Stream error")
}

// processOpenAIStream processes the streaming response from OpenAI
func processOpenAIStream(stream *openai.ChatCompletionStream, ctx context.Context) (string, error) {
	logger.LogBasic("Processing OpenAI stream")
	var fullResponse strings.Builder

	for {
		// Check if context is cancelled before attempting to read
		select {
		case <-ctx.Done():
			partialResponse := fullResponse.String()
			return "", handleOpenAIStreamRecvError(ctx.Err(), partialResponse)
		default:
			// Continue processing
		}

		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logger.LogBasic("Stream completed with EOF")
				break
			}
			// Get partial response before returning error
			partialResponse := fullResponse.String()
			return "", handleOpenAIStreamRecvError(err, partialResponse)
		}

		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			content := response.Choices[0].Delta.Content
			logger.LogVerbose("Content chunk", "text", content)
			fmt.Print(content)
			fullResponse.WriteString(content)
		}
	}

	fmt.Println()
	response := fullResponse.String()
	if response == "" {
		return "", errors.WrapError(fmt.Errorf("no content in response"), "Empty response from OpenAI")
	}

	return response, nil
}
