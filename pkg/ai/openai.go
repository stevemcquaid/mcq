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

// generateUserStoryOpenAI calls the OpenAI API with streaming
func generateUserStoryOpenAI(apiKey, featureRequest, modelID string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting OpenAI API request", "model", modelID)
	modelName := getModelDisplayName(modelID)
	showConnectionProgress("OpenAI", modelName)

	config := GetUserStoryPromptConfig(featureRequest, repoContext)
	prompt := GeneratePrompt(config)
	client := openai.NewClient(apiKey)
	req := createOpenAIRequest(modelID, prompt)

	// Create a context with timeout for stream creation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		// Check for token limit errors
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "context length") || strings.Contains(errStr, "token") || strings.Contains(errStr, "maximum context") {
			logger.LogError("Token/context limit error detected", err)
			fmt.Printf("\n⚠️  Error: Context may be too large for the model\n")
			fmt.Printf("💡 Try reducing context with --no-context or specific context flags\n")
		}
		return "", errors.WrapError(err, "Failed to create streaming request")
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			logger.LogError("close stream", closeErr)
		}
	}()

	logger.LogBasic("OpenAI stream created successfully, processing")
	showStreamingProgress()
	return processOpenAIStream(stream)
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

// processOpenAIStream processes the streaming response from OpenAI
func processOpenAIStream(stream *openai.ChatCompletionStream) (string, error) {
	logger.LogBasic("Processing OpenAI stream")
	var fullResponse strings.Builder

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logger.LogBasic("Stream completed with EOF")
				break
			}
			return "", errors.WrapError(err, "Stream error")
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
