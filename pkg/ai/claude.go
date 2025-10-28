package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/logger"
)

// generateDescriptionFromTitleClaude generates a description from a title using Claude
func generateDescriptionFromTitleClaude(apiKey, title string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting Claude API request for description generation from title")
	showConnectionProgress("Anthropic", "Claude Sonnet 4.5")

	config := GetDescriptionFromTitlePromptConfig(title, repoContext)
	prompt := GeneratePrompt(config)
	request := createClaudeRequest(prompt)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", errors.WrapError(err, "failed to marshal request")
	}

	req, err := createClaudeHTTPRequest(apiKey, jsonData)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.WrapError(err, "failed to make request")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.LogError("close response body", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.WrapError(err, fmt.Sprintf("api request failed with status %d", resp.StatusCode))
		}

		// Check for token/context limit errors
		bodyStr := strings.ToLower(string(body))
		if strings.Contains(bodyStr, "context_length") || strings.Contains(bodyStr, "input too long") || strings.Contains(bodyStr, "maximum context") {
			logger.LogError("token/context limit error detected", fmt.Errorf("claude API error: %s", string(body)))
			fmt.Printf("\n⚠️  Error: Context may be too large for Claude model\n")
			fmt.Printf("💡 Try reducing context with --no-context or specific context flags\n")
		}

		return "", errors.WrapError(fmt.Errorf("api request failed with status %d: %s", resp.StatusCode, string(body)), "claude API request failed")
	}

	logger.LogBasic("Claude API request successful, processing stream")
	showStreamingProgress()
	return processClaudeStream(resp.Body)
}

// generateImprovedDescriptionClaude generates an improved description using Claude
func generateImprovedDescriptionClaude(apiKey, originalDescription string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting Claude API request for description improvement")
	showConnectionProgress("Anthropic", "Claude Sonnet 4.5")

	config := GetDescriptionImprovementPromptConfig(originalDescription, repoContext)
	prompt := GeneratePrompt(config)
	request := createClaudeRequest(prompt)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", errors.WrapError(err, "failed to marshal request")
	}

	req, err := createClaudeHTTPRequest(apiKey, jsonData)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.WrapError(err, "failed to make request")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.LogError("close response body", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.WrapError(err, fmt.Sprintf("api request failed with status %d", resp.StatusCode))
		}

		// Check for token/context limit errors
		bodyStr := strings.ToLower(string(body))
		if strings.Contains(bodyStr, "context_length") || strings.Contains(bodyStr, "input too long") || strings.Contains(bodyStr, "maximum context") {
			logger.LogError("token/context limit error detected", fmt.Errorf("claude API error: %s", string(body)))
			fmt.Printf("\n⚠️  Error: Context may be too large for Claude model\n")
			fmt.Printf("💡 Try reducing context with --no-context or specific context flags\n")
		}

		return "", errors.WrapError(fmt.Errorf("api request failed with status %d: %s", resp.StatusCode, string(body)), "claude API request failed")
	}

	logger.LogBasic("Claude API request successful, processing stream")
	showStreamingProgress()
	return processClaudeStream(resp.Body)
}

// generateUserStoryClaude calls the Anthropic API with streaming
func generateUserStoryClaude(apiKey, featureRequest string, repoContext *RepoContext) (string, error) {
	logger.LogBasic("Starting Claude API request")
	showConnectionProgress("Anthropic", "Claude Sonnet 4.5")

	config := GetUserStoryPromptConfig(featureRequest, repoContext)
	prompt := GeneratePrompt(config)
	request := createClaudeRequest(prompt)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", errors.WrapError(err, "failed to marshal request")
	}

	req, err := createClaudeHTTPRequest(apiKey, jsonData)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.WrapError(err, "failed to make request")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.LogError("close response body", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.WrapError(err, fmt.Sprintf("api request failed with status %d", resp.StatusCode))
		}

		// Check for token/context limit errors
		bodyStr := strings.ToLower(string(body))
		if strings.Contains(bodyStr, "context_length") || strings.Contains(bodyStr, "input too long") || strings.Contains(bodyStr, "maximum context") {
			logger.LogError("token/context limit error detected", fmt.Errorf("claude API error: %s", string(body)))
			fmt.Printf("\n⚠️  Error: Context may be too large for Claude model\n")
			fmt.Printf("💡 Try reducing context with --no-context or specific context flags\n")
		}

		return "", errors.WrapError(fmt.Errorf("api request failed with status %d: %s", resp.StatusCode, string(body)), "claude API request failed")
	}

	logger.LogBasic("Claude API request successful, processing stream")
	showStreamingProgress()
	return processClaudeStream(resp.Body)
}

// createClaudeRequest creates the Anthropic API request
func createClaudeRequest(prompt string) AnthropicRequest {
	return AnthropicRequest{
		Model:     "claude-sonnet-4-5-20250929",
		MaxTokens: 4000,
		Stream:    true,
		Messages: []AnthropicMessage{
			{Role: "user", Content: prompt},
		},
	}
}

// createClaudeHTTPRequest creates the HTTP request for Claude API
func createClaudeHTTPRequest(apiKey string, jsonData []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", AnthropicAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.WrapError(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", AnthropicVersion)

	return req, nil
}

// processClaudeStream processes the streaming response from Claude
func processClaudeStream(body io.ReadCloser) (string, error) {
	logger.LogBasic("Processing Claude stream")
	var fullResponse strings.Builder
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()
		logger.LogVerbose("Stream line", "line", line)

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			logger.LogVerbose("Data event", "data", data)

			if data == "[DONE]" {
				logger.LogBasic("Stream completed with [DONE] marker")
				break
			}

			var event StreamingEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				logger.LogError("JSON unmarshal", err)
				continue
			}

			if event.Type == "content_block_delta" && event.Delta.Text != "" {
				logger.LogVerbose("Content delta", "text", event.Delta.Text)
				fmt.Print(event.Delta.Text)
				fullResponse.WriteString(event.Delta.Text)
			}
		}
	}

	fmt.Println()
	if err := scanner.Err(); err != nil {
		return "", errors.WrapError(err, "error reading stream")
	}

	response := fullResponse.String()
	if response == "" {
		return "", errors.WrapError(fmt.Errorf("no content in response"), "empty response from claude")
	}

	return response, nil
}
