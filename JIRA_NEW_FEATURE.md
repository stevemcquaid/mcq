# Jira New Command

The `mcq jira new` command creates a new Jira issue from a vague user story using AI.

## Setup

1. **Set up Jira environment variables:**
   ```bash
   export JIRA_INSTANCE_URL="https://your-jira-instance.com"
   export JIRA_USERNAME="your-username"
   export JIRA_API_TOKEN="your-api-token"  # or JIRA_PASSWORD
   export JIRA_PROJECT_PREFIX="PROJ"  # Required for issue creation
   ```

2. **Set up AI environment variables (same as `mcq ai jira`):**
   ```bash
   export ANTHROPIC_API_KEY="your-anthropic-api-key"
   # or
   export OPENAI_API_KEY="your-openai-api-key"
   ```

## Usage

The command supports two syntaxes for providing the user story:

### Syntax Options

1. **Quoted Syntax**: Wrap the user story in quotes
   ```bash
   mcq jira new "Add dark mode to the application"
   ```

2. **Unquoted Syntax**: Use `--` to separate flags from the user story
   ```bash
   mcq jira new -- Add dark mode to the application
   ```

The `--` syntax is particularly useful when the user story contains special characters or when you want to avoid shell escaping issues.

### Basic Usage (Quoted)
```bash
mcq jira new "Add dark mode to the application"
```

### Basic Usage (Unquoted with --)
```bash
mcq jira new -- Add dark mode to the application
```

### With AI Model Selection (Quoted)
```bash
mcq jira new --model claude "Add user authentication"
mcq jira new --model gpt-5 "Create a dashboard for analytics"
```

### With AI Model Selection (Unquoted)
```bash
mcq jira new --model claude -- Add user authentication
mcq jira new --model gpt-5 -- Create a dashboard for analytics
```

### With Context Gathering (Quoted)
```bash
mcq jira new --auto-context "Improve the login process"
mcq jira new --include-readme --include-go-mod "Add new feature"
```

### With Context Gathering (Unquoted)
```bash
mcq jira new --auto-context -- Improve the login process
mcq jira new --include-readme --include-go-mod -- Add new feature
```

### With Verbosity (Both Syntaxes)
```bash
# Quoted
mcq jira new --verbosity 2 "Add dark mode to the application"

# Unquoted (note: use -v 8 with space, not -v8)
mcq jira new -v 8 -- install/upgrade via homebrew
mcq jira new --verbosity 8 -- install/upgrade via homebrew
```

### Complex User Stories (Unquoted)
```bash
# User story with special characters and parentheses
mcq jira new -- User should be able to install/upgrade via homebrew (might get around code signing issues)

# Multi-word description with special characters
mcq jira new --model claude -- Create API endpoint for user@email.com authentication
```

## How It Works

1. **AI Generation**: Uses the same AI functionality as `mcq ai jira` to generate a detailed user story
2. **Content Display**: Shows the generated user story on screen
3. **Clipboard Copy**: Copies the generated content to clipboard
4. **Confirmation**: Asks for confirmation before creating the Jira issue
5. **Format Conversion**: Automatically converts markdown-style formatting to Jira markup:
   - Bullet points (`-`) → Jira bullets (`*`)
   - Indented bullets (`  -`) → Jira sub-bullets (`**`)
   - Headings (`##`) → Jira headings (`h2.`)
   - Bold text (`**text**`) → Jira bold (`*text*`)
   - Code blocks (```` ``` ````) → Jira code blocks (`{code}`)
6. **Issue Creation**: Creates a new Jira issue with:
   - **Project**: Uses `JIRA_PROJECT_PREFIX` environment variable
   - **Type**: Defaults to "Story"
   - **Summary**: Extracted from the user story (the "I want..." part)
   - **Description**: Properly formatted user story content
7. **Success**: Displays the created issue key and link

## Example Output

```bash
$ mcq jira new "Add dark mode to the application"

🤖 Generating user story with Claude Sonnet 4.5...
📝 Feature request: Add dark mode to the application

💭 As a user, I want to be able to switch between light and dark themes so that I can use the application comfortably in different lighting conditions...

📋 User story copied to clipboard!

============================================================
Generated User Story:
============================================================
As a user, I want to be able to switch between light and dark themes so that I can use the application comfortably in different lighting conditions.

Acceptance Criteria:
- User can toggle between light and dark themes
- Theme preference is saved and persists across sessions
- All UI components adapt to the selected theme
- Theme switch is accessible from the main navigation

Technical Considerations:
- Implement CSS variables for theme colors
- Add theme context provider in React
- Store theme preference in localStorage
- Ensure all components support both themes
============================================================

Create Jira issue with this content? [Y/n]: y

🔧 Creating Jira issue...
✅ Jira issue created successfully: PROJ-123
🔗 You can view it at: https://your-jira-instance.com/browse/PROJ-123
```

## Formatting Conversion

The command automatically converts markdown-style formatting to Jira markup:

**Input (AI Generated):**
```
## Acceptance Criteria
- User can toggle between light and dark themes
- Theme preference is saved and persists across sessions
- All UI components adapt to the selected theme
  - Navigation bar changes color
  - Sidebar adapts to theme
- Theme switch is accessible from the main navigation

**Priority**: High
*Estimated effort*: 3-5 days
```

**Output (Jira Markup):**
```
h2. Acceptance Criteria
* User can toggle between light and dark themes
* Theme preference is saved and persists across sessions
* All UI components adapt to the selected theme
** Navigation bar changes color
** Sidebar adapts to theme
* Theme switch is accessible from the main navigation

*Priority*: High
_Estimated effort_: 3-5 days
```

## Flags

All flags from `mcq ai jira` are supported:

- `--model, -m`: AI model to use (claude, gpt-4o, gpt-5, gpt-5-mini, gpt-5-nano)
- `--verbosity, -v`: Verbosity level (0=off, 1=basic, 2=detailed, 3=verbose)
- `--auto-context`: Automatically detect and include repository context
- `--include-readme`: Include README content in context
- `--include-go-mod`: Include go.mod information in context
- `--include-commits`: Include recent commit messages in context
- `--include-structure`: Include directory structure in context
- `--include-configs`: Include configuration files in context
- `--max-commits`: Maximum number of recent commits to include
- `--no-context`: Skip context gathering entirely

## Help

To see all available options and flags:

```bash
mcq jira new --help
```

This will display all available flags including AI model selection, context gathering options, and verbosity levels.

## Requirements

- macOS (for clipboard functionality)
- Valid Jira API credentials
- Valid AI API key (Anthropic or OpenAI)
- `JIRA_PROJECT_PREFIX` environment variable must be set

## Error Handling

The command handles various error scenarios:
- Missing Jira configuration
- Missing AI API keys
- Jira API failures
- AI generation failures
- Clipboard copy failures (non-fatal)

## Integration

This command integrates seamlessly with the existing `mcq` tool:
- Uses the same Jira authentication as `mcq jira show`
- Uses the same AI functionality as `mcq ai jira`
- Follows the same command structure and flag patterns
- Maintains consistency with existing error handling and logging
