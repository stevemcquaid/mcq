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
	"path/filepath"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// ============================================================================
// TYPES AND CONSTANTS
// ============================================================================

// ModelConfig represents configuration for an AI model
type ModelConfig struct {
	Name        string
	Provider    string
	APIKey      string
	ModelID     string
	Description string
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Logger *slog.Logger
}

// ContextConfig holds configuration for context gathering
type ContextConfig struct {
	AutoDetect       bool
	IncludeReadme    bool
	IncludeGoMod     bool
	IncludeCommits   bool
	IncludeStructure bool
	IncludeConfigs   bool
	MaxCommits       int
	MaxFileSize      int64 // in bytes
}

// RepoContext holds gathered repository context
type RepoContext struct {
	ProjectName        string
	GoVersion          string
	ModulePath         string
	Dependencies       []string
	Readme             string
	RecentCommits      []string
	DirectoryStructure string
	ConfigFiles        map[string]string
	ProjectType        string
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

// Constants
const (
	DefaultMaxTokens = 4000
	GPT5Prefix       = "gpt-5"
	AnthropicAPIURL  = "https://api.anthropic.com/v1/messages"
	AnthropicVersion = "2023-06-01"
)

// Log levels
var (
	LogLevelOff      = slog.LevelError - 1
	LogLevelBasic    = slog.LevelInfo
	LogLevelDetailed = slog.LevelDebug
	LogLevelVerbose  = slog.LevelDebug - 1
)

// Global logger configuration
var loggerConfig LoggerConfig

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

// ============================================================================
// LOGGING FUNCTIONS
// ============================================================================

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

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
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

// logVerbose logs verbose information
func logVerbose(msg string, args ...interface{}) {
	if loggerConfig.Logger != nil {
		loggerConfig.Logger.Log(context.Background(), LogLevelVerbose, msg, args...)
	}
}

// logError logs error information
func logError(operation string, err error) {
	if loggerConfig.Logger != nil {
		loggerConfig.Logger.Error("Operation failed", "operation", operation, "error", err)
	}
}

// ============================================================================
// MAIN AI JIRA FUNCTION
// ============================================================================

// AIJira converts a vague feature request to a user story and copies it to clipboard
func AIJira(args []string, modelFlag string, verbosityLevel int, contextConfig ContextConfig) error {
	setupLogger(verbosityLevel)

	featureRequest := strings.Join(args, " ")
	logBasic("Starting AIJira", "feature_request", featureRequest)

	// Gather repository context
	repoContext := gatherContextIfNeeded(contextConfig)

	// Select and configure model
	selectedModel, err := selectModel(modelFlag)
	if err != nil {
		return err
	}

	// Generate user story
	userStory, err := generateUserStory(selectedModel, featureRequest, repoContext)
	if err != nil {
		return fmt.Errorf("failed to generate user story: %w", err)
	}

	// Display and copy result
	return displayAndCopyResult(userStory)
}

// ============================================================================
// CONTEXT GATHERING
// ============================================================================

// gatherContextIfNeeded gathers repository context if any context options are enabled
func gatherContextIfNeeded(config ContextConfig) *RepoContext {
	if !shouldGatherContext(config) {
		return nil
	}

	logBasic("Gathering repository context")
	repoContext, err := gatherRepoContext(config)
	if err != nil {
		logError("context gathering", err)
		return nil
	}

	logBasic("Repository context gathered successfully")
	return repoContext
}

// shouldGatherContext determines if any context should be gathered
func shouldGatherContext(config ContextConfig) bool {
	return config.AutoDetect || config.IncludeReadme || config.IncludeGoMod ||
		config.IncludeCommits || config.IncludeStructure || config.IncludeConfigs
}

// gatherRepoContext gathers repository context based on configuration
func gatherRepoContext(config ContextConfig) (*RepoContext, error) {
	ctx := &RepoContext{
		ConfigFiles: make(map[string]string),
	}

	// Apply auto-detect settings if enabled
	config = applyAutoDetectSettings(config)

	// Gather all context components
	gatherContextComponents(ctx, config)

	// Determine project type
	ctx.ProjectType = determineProjectType(ctx)

	return ctx, nil
}

// applyAutoDetectSettings applies auto-detect settings to the config
func applyAutoDetectSettings(config ContextConfig) ContextConfig {
	if config.AutoDetect {
		config.IncludeReadme = true
		config.IncludeGoMod = true
		config.IncludeCommits = true
		config.IncludeStructure = true
		config.IncludeConfigs = true
		config.MaxCommits = 10
		config.MaxFileSize = 50 * 1024 // 50KB
	}
	return config
}

// gatherContextComponents gathers all enabled context components
func gatherContextComponents(ctx *RepoContext, config ContextConfig) {
	gatherComponent(ctx, config.IncludeGoMod, "Go module info", func() error {
		return gatherGoModuleInfo(ctx)
	})

	gatherComponent(ctx, config.IncludeReadme, "README", func() error {
		return gatherReadme(ctx)
	})

	gatherComponent(ctx, config.IncludeCommits, "recent commits", func() error {
		return gatherRecentCommits(ctx, config.MaxCommits)
	})

	gatherComponent(ctx, config.IncludeStructure, "directory structure", func() error {
		return gatherDirectoryStructure(ctx)
	})

	gatherComponent(ctx, config.IncludeConfigs, "config files", func() error {
		return gatherConfigFiles(ctx, config.MaxFileSize)
	})
}

// gatherComponent is a helper to gather a context component with error logging
func gatherComponent(_ *RepoContext, shouldGather bool, componentName string, gatherFunc func() error) {
	if !shouldGather {
		return
	}

	if err := gatherFunc(); err != nil {
		if loggerConfig.Logger != nil {
			loggerConfig.Logger.Debug("Failed to gather "+componentName, "error", err)
		}
	}
}

// gatherGoModuleInfo extracts information from go.mod
func gatherGoModuleInfo(ctx *RepoContext) error {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("go.mod not found: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Extract module path
		if strings.HasPrefix(line, "module ") {
			ctx.ModulePath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			ctx.ProjectName = filepath.Base(ctx.ModulePath)
		}

		// Extract Go version
		if strings.HasPrefix(line, "go ") {
			ctx.GoVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}

		// Extract dependencies
		if isDependencyLine(line) {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				ctx.Dependencies = append(ctx.Dependencies, parts[0])
			}
		}
	}

	return nil
}

// isDependencyLine checks if a line contains a dependency
func isDependencyLine(line string) bool {
	return strings.HasPrefix(line, "require ") ||
		(strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "\t//") &&
			!strings.HasPrefix(line, "go ") && !strings.HasPrefix(line, "module "))
}

// gatherReadme extracts README content
func gatherReadme(ctx *RepoContext) error {
	readmeFiles := []string{"README.md", "README.rst", "README.txt", "README"}

	for _, filename := range readmeFiles {
		content, err := os.ReadFile(filename)
		if err == nil {
			ctx.Readme = string(content)
			return nil
		}
	}

	return fmt.Errorf("no README file found")
}

// gatherRecentCommits gets recent commit messages
func gatherRecentCommits(ctx *RepoContext, maxCommits int) error {
	cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("-n%d", maxCommits))
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line != "" {
			ctx.RecentCommits = append(ctx.RecentCommits, line)
		}
	}

	return nil
}

// gatherDirectoryStructure gets a high-level directory structure
func gatherDirectoryStructure(ctx *RepoContext) error {
	var structure strings.Builder

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if shouldSkipPath(path) {
			return getSkipAction()
		}

		addToStructure(&structure, path, info)
		return nil
	})

	ctx.DirectoryStructure = structure.String()
	return err
}

// shouldSkipPath determines if a path should be skipped
func shouldSkipPath(path string) bool {
	// Skip hidden directories
	if strings.HasPrefix(path, ".") && path != "." {
		return true
	}

	// Skip common directories that don't add value
	skipDirs := []string{"vendor", "node_modules", ".git", "build", "dist", "target", "bin", "obj"}
	for _, skipDir := range skipDirs {
		if strings.Contains(path, skipDir) {
			return true
		}
	}

	return false
}

// getSkipAction returns the appropriate skip action for a directory
func getSkipAction() error {
	return filepath.SkipDir
}

// addToStructure adds a path to the directory structure
func addToStructure(structure *strings.Builder, path string, info os.FileInfo) {
	depth := strings.Count(path, string(filepath.Separator))
	indent := strings.Repeat("  ", depth)

	if info.IsDir() {
		if depth < 3 { // Limit depth to avoid too much detail
			fmt.Fprintf(structure, "%s%s/\n", indent, info.Name())
		}
	} else if isImportantFile(path) {
		fmt.Fprintf(structure, "%s%s\n", indent, info.Name())
	}
}

// isImportantFile determines if a file is important for context
func isImportantFile(path string) bool {
	importantExts := []string{".go", ".md", ".yaml", ".yml", ".json", ".toml", ".env", ".dockerfile", "Dockerfile", "Makefile"}
	importantNames := []string{"go.mod", "go.sum", "README", "LICENSE", "CHANGELOG", "Dockerfile", "Makefile", ".gitignore"}

	ext := filepath.Ext(path)
	for _, importantExt := range importantExts {
		if ext == importantExt {
			return true
		}
	}

	filename := filepath.Base(path)
	for _, importantName := range importantNames {
		if filename == importantName {
			return true
		}
	}

	return false
}

// gatherConfigFiles collects relevant configuration files
func gatherConfigFiles(ctx *RepoContext, maxSize int64) error {
	configFiles := []string{
		"go.mod", "go.sum", "Makefile", "Dockerfile", ".dockerignore",
		"docker-compose.yml", "docker-compose.yaml", ".env", ".env.example",
		"config.yaml", "config.yml", "config.json", ".gitignore",
	}

	for _, filename := range configFiles {
		if info, err := os.Stat(filename); err == nil && info.Size() <= maxSize {
			if content, err := os.ReadFile(filename); err == nil {
				ctx.ConfigFiles[filename] = string(content)
			}
		}
	}

	return nil
}

// determineProjectType analyzes the repository to determine project type
func determineProjectType(ctx *RepoContext) string {
	// Check README for project type indicators
	if strings.Contains(ctx.Readme, "CLI") || strings.Contains(ctx.Readme, "command") {
		return "CLI Tool"
	}
	if strings.Contains(ctx.Readme, "API") || strings.Contains(ctx.Readme, "server") {
		return "Web API"
	}
	if strings.Contains(ctx.Readme, "library") || strings.Contains(ctx.Readme, "package") {
		return "Library"
	}

	// Check dependencies for clues
	for _, dep := range ctx.Dependencies {
		if strings.Contains(dep, "gin") || strings.Contains(dep, "echo") || strings.Contains(dep, "fiber") {
			return "Web API"
		}
		if strings.Contains(dep, "cobra") || strings.Contains(dep, "cli") {
			return "CLI Tool"
		}
	}

	// Check directory structure
	if strings.Contains(ctx.DirectoryStructure, "cmd/") {
		return "CLI Tool"
	}
	if strings.Contains(ctx.DirectoryStructure, "api/") || strings.Contains(ctx.DirectoryStructure, "server/") {
		return "Web API"
	}

	return "Go Application"
}

// formatContextForPrompt formats the repository context for inclusion in AI prompts
func formatContextForPrompt(ctx *RepoContext) string {
	if ctx == nil {
		return ""
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("\n## Repository Context\n\n")

	// Project information
	contextBuilder.WriteString("### Project Information\n")
	contextBuilder.WriteString(fmt.Sprintf("- **Project Name**: %s\n", ctx.ProjectName))
	contextBuilder.WriteString(fmt.Sprintf("- **Module Path**: %s\n", ctx.ModulePath))
	contextBuilder.WriteString(fmt.Sprintf("- **Go Version**: %s\n", ctx.GoVersion))
	contextBuilder.WriteString(fmt.Sprintf("- **Project Type**: %s\n\n", ctx.ProjectType))

	// Dependencies
	if len(ctx.Dependencies) > 0 {
		contextBuilder.WriteString("### Key Dependencies\n")
		for _, dep := range ctx.Dependencies[:minInt(10, len(ctx.Dependencies))] {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", dep))
		}
		contextBuilder.WriteString("\n")
	}

	// README excerpt
	if ctx.Readme != "" {
		contextBuilder.WriteString("### Project Overview\n")
		readmeExcerpt := ctx.Readme
		if len(readmeExcerpt) > 1000 {
			readmeExcerpt = readmeExcerpt[:1000] + "..."
		}
		contextBuilder.WriteString(readmeExcerpt)
		contextBuilder.WriteString("\n\n")
	}

	// Recent commits
	if len(ctx.RecentCommits) > 0 {
		contextBuilder.WriteString("### Recent Development Activity\n")
		for _, commit := range ctx.RecentCommits[:minInt(5, len(ctx.RecentCommits))] {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", commit))
		}
		contextBuilder.WriteString("\n")
	}

	// Directory structure
	if ctx.DirectoryStructure != "" {
		contextBuilder.WriteString("### Project Structure\n")
		contextBuilder.WriteString("```\n")
		contextBuilder.WriteString(ctx.DirectoryStructure)
		contextBuilder.WriteString("\n```\n\n")
	}

	// Configuration files
	if len(ctx.ConfigFiles) > 0 {
		contextBuilder.WriteString("### Configuration Files\n")
		for filename, content := range ctx.ConfigFiles {
			contextBuilder.WriteString(fmt.Sprintf("**%s**:\n", filename))
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			contextBuilder.WriteString("```\n")
			contextBuilder.WriteString(content)
			contextBuilder.WriteString("\n```\n\n")
		}
	}

	return contextBuilder.String()
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PromptForContext interactively asks the user if they want to include context
func PromptForContext() ContextConfig {
	fmt.Println("\n🔍 Would you like to include repository context to improve the user story?")
	fmt.Println()
	fmt.Println("This will include:")
	fmt.Println("  📄 README content and project description")
	fmt.Println("  🔧 Go module information and dependencies")
	fmt.Println("  📝 Recent commit messages (last 10 commits)")
	fmt.Println("  📁 Directory structure overview")
	fmt.Println("  ⚙️  Configuration files (go.mod, Makefile, etc.)")
	fmt.Println()
	fmt.Println("This helps generate more accurate and contextually relevant user stories.")
	fmt.Println()

	// Ask if they want context
	fmt.Print("Include repository context? (Y/n): ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println("\nError reading input, skipping context gathering.")
		return ContextConfig{}
	}

	if strings.ToLower(response) == "n" || strings.ToLower(response) == "no" {
		fmt.Println("Skipping context gathering.")
		return ContextConfig{}
	}

	fmt.Println("✅ Including repository context...")
	return ContextConfig{
		AutoDetect:       true,
		IncludeReadme:    true,
		IncludeGoMod:     true,
		IncludeCommits:   true,
		IncludeStructure: true,
		IncludeConfigs:   true,
		MaxCommits:       10,
		MaxFileSize:      50 * 1024, // 50KB default
	}
}

// ============================================================================
// MODEL SELECTION
// ============================================================================

// selectModel determines which AI model to use
func selectModel(modelFlag string) (ModelConfig, error) {
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")

	logDetailed("API Keys",
		"anthropic", maskAPIKey(anthropicAPIKey),
		"openai", maskAPIKey(openaiAPIKey))

	// If model is explicitly specified, validate and return it
	if modelFlag != "" {
		return selectExplicitModel(modelFlag, anthropicAPIKey, openaiAPIKey)
	}

	// Auto-detect based on available API keys
	return selectModelByAvailability(anthropicAPIKey, openaiAPIKey)
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
			return ModelConfig{}, fmt.Errorf("ANTHROPIC_API_KEY is required for Claude model")
		}
		model.APIKey = anthropicAPIKey
	case "openai":
		if openaiAPIKey == "" {
			return ModelConfig{}, fmt.Errorf("OPENAI_API_KEY is required for %s model", model.Name)
		}
		model.APIKey = openaiAPIKey
	}

	logBasic("Selected model", "name", model.Name, "provider", model.Provider)
	return model, nil
}

// selectModelByAvailability selects a model based on available API keys
func selectModelByAvailability(anthropicAPIKey, openaiAPIKey string) (ModelConfig, error) {
	hasAnthropic := anthropicAPIKey != ""
	hasOpenAI := openaiAPIKey != ""

	if !hasAnthropic && !hasOpenAI {
		return ModelConfig{}, fmt.Errorf("no API keys found. Please set either ANTHROPIC_API_KEY or OPENAI_API_KEY")
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
		return ModelConfig{}, fmt.Errorf("invalid input: %w", err)
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

// ============================================================================
// USER STORY GENERATION
// ============================================================================

// generateUserStory generates a user story using the specified model
func generateUserStory(model ModelConfig, featureRequest string, repoContext *RepoContext) (string, error) {
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
	basePrompt := `Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations.

Feature Request: %s

Please provide a comprehensive user story with:
1. The main user story in the specified format
2. Acceptance criteria
3. Any relevant technical notes or considerations`

	// Add repository context if available
	if repoContext != nil {
		contextInfo := formatContextForPrompt(repoContext)
		basePrompt += contextInfo
	}

	return fmt.Sprintf(basePrompt, featureRequest)
}

// ============================================================================
// CLAUDE API INTEGRATION
// ============================================================================

// generateUserStoryClaude calls the Anthropic API with streaming
func generateUserStoryClaude(apiKey, featureRequest string, repoContext *RepoContext) (string, error) {
	logBasic("Starting Claude API request")
	showConnectionProgress("Anthropic", "Claude Sonnet 4.5")

	prompt := createPrompt(featureRequest, repoContext)
	request := createClaudeRequest(prompt)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := createClaudeHTTPRequest(apiKey, jsonData)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logError("close response body", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("API request failed with status %d: failed to read response body: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	logBasic("Claude API request successful, processing stream")
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", AnthropicVersion)

	return req, nil
}

// processClaudeStream processes the streaming response from Claude
func processClaudeStream(body io.ReadCloser) (string, error) {
	logBasic("Processing Claude stream")
	var fullResponse strings.Builder
	scanner := bufio.NewScanner(body)

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

			if event.Type == "content_block_delta" && event.Delta.Text != "" {
				logVerbose("Content delta", "text", event.Delta.Text)
				fmt.Print(event.Delta.Text)
				fullResponse.WriteString(event.Delta.Text)
			}
		}
	}

	fmt.Println()
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading stream: %w", err)
	}

	response := fullResponse.String()
	if response == "" {
		return "", fmt.Errorf("no content in response")
	}

	return response, nil
}

// ============================================================================
// OPENAI API INTEGRATION
// ============================================================================

// generateUserStoryOpenAI calls the OpenAI API with streaming
func generateUserStoryOpenAI(apiKey, featureRequest, modelID string, repoContext *RepoContext) (string, error) {
	logBasic("Starting OpenAI API request", "model", modelID)
	modelName := getModelDisplayName(modelID)
	showConnectionProgress("OpenAI", modelName)

	prompt := createPrompt(featureRequest, repoContext)
	client := openai.NewClient(apiKey)
	req := createOpenAIRequest(modelID, prompt)

	// Create a context with timeout for stream creation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create streaming request: %w", err)
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			logError("close stream", closeErr)
		}
	}()

	logBasic("OpenAI stream created successfully, processing")
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
	logBasic("Processing OpenAI stream")
	var fullResponse strings.Builder

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logBasic("Stream completed with EOF")
				break
			}
			return "", fmt.Errorf("stream error: %w", err)
		}

		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			content := response.Choices[0].Delta.Content
			logVerbose("Content chunk", "text", content)
			fmt.Print(content)
			fullResponse.WriteString(content)
		}
	}

	fmt.Println()
	response := fullResponse.String()
	if response == "" {
		return "", fmt.Errorf("no content in response")
	}

	return response, nil
}

// ============================================================================
// DISPLAY AND CLIPBOARD FUNCTIONS
// ============================================================================

// displayAndCopyResult copies the user story to clipboard and displays it
func displayAndCopyResult(userStory string) error {
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

// showConnectionProgress displays progress during API connection setup
func showConnectionProgress(provider, modelName string) {
	fmt.Printf("🔌 Connecting to %s API (%s)...\n", provider, modelName)
}

// showStreamingProgress displays progress when streaming starts
func showStreamingProgress() {
	fmt.Print("💭 ")
}
