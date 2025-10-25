package ai

// ModelConfig represents configuration for an AI model
type ModelConfig struct {
	Name        string
	Provider    string
	APIKey      string
	ModelID     string
	Description string
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

// StreamingEvent represents a streaming response event from Anthropic
type StreamingEvent struct {
	Type         string `json:"type"`
	ContentBlock struct {
		Text string `json:"text"`
	} `json:"content_block"`
	Delta struct {
		Text string `json:"text"`
	} `json:"delta"`
}

// Constants
const (
	DefaultMaxTokens = 4000
	GPT5Prefix       = "gpt-5"
	AnthropicAPIURL  = "https://api.anthropic.com/v1/messages"
	AnthropicVersion = "2023-06-01"
)
