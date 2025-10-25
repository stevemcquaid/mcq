# AI Jira Command

The `mcq ai jira` command converts vague feature requests into detailed user stories using AI models and copies the result to your clipboard.

## Setup

1. **Get API Keys**: 
   - **Claude (Anthropic)**: Sign up at [Anthropic Console](https://console.anthropic.com/) and get your API key.
   - **OpenAI Models**: Sign up at [OpenAI Platform](https://platform.openai.com/) and get your API key for GPT-4o, GPT-5, GPT-5 Mini, or GPT-5 Nano.

2. **Set Environment Variables**: Set your API keys as environment variables:
   ```bash
   # For Claude
   export ANTHROPIC_API_KEY='your-anthropic-api-key-here'
   
   # For OpenAI models (GPT-4o, GPT-5, GPT-5 Mini, GPT-5 Nano)
   export OPENAI_API_KEY='your-openai-api-key-here'
   
   # Or set both to choose between models
   ```

   To make them persistent, add these lines to your shell profile (e.g., `~/.zshrc` or `~/.bash_profile`).

3. **Customize Prompts (Optional)**: You can customize the AI prompts using Go templates:
   ```bash
   # Generate example template files
   mcq templates generate ./my-templates
   
   # Set custom template directory
   export MCQ_PROMPTS_DIR='./my-templates'
   
   # Validate templates
   mcq templates validate
   ```

## Usage

```bash
mcq ai jira <vague feature request>
```

### Model Selection

**Auto-detection (Recommended):**
```bash
# Automatically uses available model(s)
mcq ai jira "Add a dark mode to the application"
```

**Explicit model selection:**
```bash
# Use Claude specifically
mcq ai jira --model claude "Add a dark mode to the application"

# Use GPT-4o specifically  
mcq ai jira --model gpt-4o "Add a dark mode to the application"

# Use GPT-5 (latest and most advanced)
mcq ai jira --model gpt-5 "Add a dark mode to the application"

# Use GPT-5 Mini (faster and more cost-effective)
mcq ai jira --model gpt-5-mini "Add a dark mode to the application"

# Use GPT-5 Nano (optimized for simple tasks)
mcq ai jira --model gpt-5-nano "Add a dark mode to the application"

# Set verbosity levels for structured logging
mcq ai jira --verbosity 1 "Add a dark mode to the application"  # Basic logging
mcq ai jira --verbosity 2 "Add a dark mode to the application"  # Detailed logging
mcq ai jira --verbosity 3 "Add a dark mode to the application"  # Verbose logging (includes streaming details)
```

**Interactive selection:**
When both API keys are available, you'll be prompted to choose:
```
🔑 Both Claude and OpenAI API keys are available.
Which model would you like to use?
1. Claude Sonnet 4.5 (Anthropic)
2. GPT-4o (OpenAI)
3. GPT-5 (OpenAI) - Full power, best for complex tasks
4. GPT-5 Mini (OpenAI) - Faster and more cost-effective
5. GPT-5 Nano (OpenAI) - Optimized for simple tasks
Enter choice (1-5):
```

### Examples

```bash
# Basic usage (auto-detects model)
mcq ai jira "Add a dark mode to the application"

# More complex request with specific model
mcq ai jira --model claude "Improve the user login process to make it more secure and user-friendly"

# Multiple words
mcq ai jira "Create a dashboard that shows real-time analytics for our users"
```

## Output

The command will:
1. Show progress indicators and status updates
2. Stream the response in real-time as Claude generates it
3. Generate a detailed user story following the format: "As a [user type], I want [goal] so that [benefit]"
4. Include additional information such as:
   - Acceptance criteria
   - Technical considerations
   - Priority level (if applicable)
5. Copy the generated user story to your clipboard
6. Display the complete user story in the terminal with formatting

### Progress Indicators

- 🤖 Shows when starting generation with the selected AI model
- 📝 Displays the feature request being processed
- ⏳ Indicates processing is in progress
- 💭 Shows real-time streaming output as it's generated
- 📋 Confirms clipboard operation
- ✅ Indicates successful completion
- 🔑 Shows model selection prompt when both API keys are available

## Structured Logging

The `--verbosity` flag provides configurable levels of structured logging using Go's standard `slog` package:

```bash
mcq ai jira --verbosity 1 "Add a dark mode to the application"  # Basic logging
mcq ai jira --verbosity 2 "Add a dark mode to the application"  # Detailed logging
mcq ai jira --verbosity 3 "Add a dark mode to the application"  # Verbose logging
```

**Verbosity Levels:**

- **Level 0 (Off)**: No logging output (default)
- **Level 1 (Basic)**: Essential process information, errors, and completion status
- **Level 2 (Detailed)**: API details, request/response information, processing summaries
- **Level 3 (Verbose)**: All details including individual streaming chunks and real-time processing

**Level 1 (Basic) includes:**
- Process start/completion messages
- Model selection results
- Error messages with structured context
- Success confirmations

**Level 2 (Detailed) includes:**
- API key validation (masked for security)
- Model selection process details
- API request/response details with structured fields
- Processing statistics
- Clipboard operation status

**Level 3 (Verbose) includes:**
- Individual stream events and data chunks
- Real-time content processing details
- Detailed API response parsing
- All streaming debug information

**Structured Logging Benefits:**
- Consistent key-value format for easy parsing
- Better error context and debugging
- Standard Go logging practices
- Easy integration with log aggregation systems

**Example structured logging output (Level 2 - Detailed):**
```
level=INFO msg="Starting AIJira" feature_request="Add a dark mode to the application"
level=INFO msg="Configuration" model_flag="" verbosity_level=2
level=DEBUG msg="API Keys" anthropic="***1234" openai="***5678"
level=INFO msg="Selecting model"
level=DEBUG msg="Auto-detection" anthropic_available=true openai_available=true
level=DEBUG msg="Both API keys available, prompting user for selection"
level=DEBUG msg="API operation" operation="Claude Request" details="Model: claude-sonnet-4-5-20250929, MaxTokens: 4000, Stream: true"
level=DEBUG msg="Stream processing complete" events=15 content_blocks=12
level=INFO msg="Successfully copied to clipboard"
```

**Example structured logging output (Level 3 - Verbose):**
```
level=INFO msg="Starting AIJira" feature_request="Add a dark mode to the application"
level=INFO msg="Configuration" model_flag="" verbosity_level=3
level=DEBUG msg="API Keys" anthropic="***1234" openai="***5678"
level=INFO msg="Selecting model"
level=DEBUG msg="Auto-detection" anthropic_available=true openai_available=true
level=DEBUG msg="Both API keys available, prompting user for selection"
level=DEBUG msg="API operation" operation="Claude Request" details="Model: claude-sonnet-4-5-20250929, MaxTokens: 4000, Stream: true"
level=DEBUG-1 msg="Stream line" line="data: {\"type\":\"content_block_delta\",\"delta\":{\"text\":\"As a user\"}}"
level=DEBUG-1 msg="Content delta" count=1 text="As a user"
level=DEBUG msg="Stream processing complete" events=15 content_blocks=12
level=INFO msg="Successfully copied to clipboard"
```

## Template Customization

The AI prompts can be customized using Go templates. This allows you to modify the behavior and output format of the AI without changing the code.

### Available Template Variables

**For User Story Templates:**
- `{{.FeatureRequest}}` - The user's feature request
- `{{.RepositoryContext}}` - Repository information (if available)
- `{{.ProjectName}}` - Project name from go.mod
- `{{.ModulePath}}` - Module path from go.mod
- `{{.GoVersion}}` - Go version from go.mod
- `{{.ProjectType}}` - Detected project type
- `{{.Readme}}` - README content
- `{{.RecentCommits}}` - Recent commit messages
- `{{.Dependencies}}` - Go dependencies
- `{{.DirectoryStructure}}` - Directory structure
- `{{.ConfigFiles}}` - Configuration files content
- `{{.Now}}` - Current timestamp

**For Title Extraction Templates:**
- `{{.FeatureRequest}}` - The original feature request
- `{{.UserStory}}` - The generated user story
- `{{.Now}}` - Current timestamp

### Template Commands

```bash
# Generate example template files
mcq templates generate ./my-templates

# Validate template syntax
mcq templates validate

# List available prompt types
mcq templates list
```

### Example Custom Template

```gotemplate
{{/* Custom User Story Template */}}
You are an expert product manager. Convert this feature request into a detailed user story.

Feature: {{.FeatureRequest}}

{{if .RepositoryContext}}
Project Context:
- Project: {{.ProjectName}}
- Type: {{.ProjectType}}
- Go Version: {{.GoVersion}}
{{end}}

Please provide:
1. A user story in "As a [user], I want [goal] so that [benefit]" format
2. Acceptance criteria
3. Technical considerations
4. Keep under 500 words

User Story:
```

## Error Handling

The command will show clear error messages for:
- Missing API keys (when no `ANTHROPIC_API_KEY` or `OPENAI_API_KEY` is set)
- Invalid model selection
- API request failures (with detailed error text when debug is enabled)
- Clipboard operation failures
- Network connectivity issues
- Invalid user choices during model selection

### Common Issues

**OpenAI API Parameter Errors:**
If you encounter errors like "Unsupported parameter: 'max_tokens' is not supported with this model", the command automatically handles this by:
- Using `max_tokens` for older models (GPT-4o)
- Omitting the parameter for newer models (GPT-5 variants)
- Providing detailed debug output to help troubleshoot

**Library Updates:**
The command uses the go-openai library. If you encounter API compatibility issues, you may need to update the library:
```bash
go get -u github.com/sashabaranov/go-openai
```

## Requirements

- macOS (for clipboard functionality using `pbcopy`)
- Valid API key for at least one model:
  - Anthropic API key for Claude Sonnet 4.5
  - OpenAI API key for GPT-4o, GPT-5, GPT-5 Mini, or GPT-5 Nano
- Internet connection

## Model Comparison

| Model | Provider | Best For | Speed | Cost |
|-------|----------|----------|-------|------|
| Claude Sonnet 4.5 | Anthropic | Complex reasoning, coding | Medium | Medium |
| GPT-4o | OpenAI | General purpose | Fast | Low |
| GPT-5 | OpenAI | Complex tasks, agentic work | Medium | High |
| GPT-5 Mini | OpenAI | Most applications | Fast | Medium |
| GPT-5 Nano | OpenAI | Simple tasks | Very Fast | Low |
