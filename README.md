# mcq
A golang development helper with AI-powered features and customizable prompt templates. Don't memorize commands when you can `mcq lint`

# Usage
`mcq help`
```
$ mcq help
This application provides shortcuts to common development tasks

Usage:
  mcq [command]

Available Commands:
  ai          AI-powered commands
  all         Run everything
  build       -> go build
  ci          Run almost everything
  clean       -> fmt deps vet
  cover       -> go tool cover
  deps        -> go mod tidy, download, vendor
  docker      docker build, run, push
  fmt         -> go fmt
  help        Help about any command
  install     -> go install
  lint        -> golangci-lint, staticcheck
  log         -> ~git log --graph --oneline --decorate --all
  run         -> go run main.go
  setup       install dependencies
  test        -> go test
  version     Version

Flags:
  -h, --help   help for mcq

Use "mcq [command] --help" for more information about a command.
```

## AI Commands

### `mcq ai jira <feature request>`

Convert vague feature requests into detailed user stories using AI models and copy to clipboard.

**Setup:**
```bash
# For Claude (Anthropic)
export ANTHROPIC_API_KEY='your-anthropic-api-key-here'

# For GPT-4o (OpenAI) 
export OPENAI_API_KEY='your-openai-api-key-here'

# Or set both to choose between models
```

**Usage:**
```bash
# Auto-detect model based on available API keys
mcq ai jira "Add a dark mode to the application"

# Specify model explicitly
mcq ai jira --model claude "Add a dark mode to the application"
mcq ai jira --model gpt-4o "Add a dark mode to the application"
mcq ai jira --model gpt-5 "Add a dark mode to the application"
mcq ai jira --model gpt-5-mini "Add a dark mode to the application"
mcq ai jira --model gpt-5-nano "Add a dark mode to the application"

# Set verbosity levels for structured logging
mcq ai jira --verbosity 1 "Add a dark mode to the application"  # Basic logging
mcq ai jira --verbosity 2 "Add a dark mode to the application"  # Detailed logging
mcq ai jira --verbosity 3 "Add a dark mode to the application"  # Verbose logging (includes streaming details)
```

**Supported Models:**
- **Claude Sonnet 4.5** (Anthropic) - Latest Claude model
- **GPT-4o** (OpenAI) - Previous generation GPT model
- **GPT-5** (OpenAI) - Latest and most advanced GPT model
- **GPT-5 Mini** (OpenAI) - Faster and more cost-effective GPT-5 variant
- **GPT-5 Nano** (OpenAI) - Optimized for simple tasks with minimal latency

This will generate a user story in the format "As a [user type], I want [goal] so that [benefit]" with additional acceptance criteria and technical considerations, then copy it to your clipboard.

For more details, see [AI_FEATURE.md](AI_FEATURE.md).

## Template Customization

Customize AI prompts using Go templates without modifying the code.

**Setup:**
```bash
# Generate example template files
mcq templates generate ./my-templates

# Set custom template directory
export MCQ_PROMPTS_DIR='./my-templates'

# Validate templates
mcq templates validate
```

**Available Commands:**
- `mcq templates generate [dir]` - Generate example template files
- `mcq templates validate` - Validate template syntax
- `mcq templates list` - List available prompt types

**Template Variables:**
- `{{.FeatureRequest}}` - User's feature request
- `{{.RepositoryContext}}` - Repository information
- `{{.ProjectName}}`, `{{.ModulePath}}`, `{{.GoVersion}}` - Project details
- `{{.Readme}}`, `{{.RecentCommits}}`, `{{.Dependencies}}` - Context data
- `{{.Now}}` - Current timestamp

For more details, see [TEMPLATES.md](TEMPLATES.md).

# TODO
* [x] Mechanism to fail fast during commands running. If error, it should quit. (OrderedRunner)
* [ ] Mechanism for pretty printing text to screen. Likely a writer library/passed around with global defaults for different types of messages
* [ ] Mechanism for parallelization of tasks than can be completed together
* [ ] Simplify colorwriter
