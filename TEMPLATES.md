# AI Prompt Templates

This document describes how to customize AI prompts using Go templates.

## Overview

The AI prompt system now supports custom templates that can be modified without changing the code. Templates use Go's `html/template` package and support variables and conditional logic.

## Environment Variable

Set the `MCQ_PROMPTS_DIR` environment variable to point to a directory containing your custom template files:

```bash
export MCQ_PROMPTS_DIR="/path/to/your/templates"
```

If not set, the system will use built-in default templates.

## Template Files

Template files should be named according to the prompt type:

- `user_story.tpl` - For user story generation
- `title_extraction.tpl` - For JIRA title extraction

## Available Variables

### For User Story Templates

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

### For Title Extraction Templates

- `{{.FeatureRequest}}` - The original feature request
- `{{.UserStory}}` - The generated user story
- `{{.Now}}` - Current timestamp

## Template Functions

- `{{formatContext .RepositoryContext}}` - Formats repository context for inclusion in prompts

## Example Templates

### User Story Template

```gotemplate
{{/* 
User Story Generation Template
Available variables listed above
*/}}
Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations. Provide ONLY the user story. 

Please provide a comprehensive user story:
1. With the main user story in the specified format
2. With acceptance criteria
3. With any relevant technical notes or considerations
4. Keep the total output under 1000 words

Do NOT add any additional questions or commentary. 
The response must ONLY be the user story. 
NOTHING ELSE.

Feature Request: {{.FeatureRequest}}
{{if .RepositoryContext}}
{{formatContext .RepositoryContext}}
{{end}}
```

### Title Extraction Template

```gotemplate
{{/* 
Title Extraction Template
Available variables listed above
*/}}
Create a NEW concise, clear title (maximum 100 characters) for a Jira issue from the following user story and old title. The new title should be action-oriented and summarize the main goal or feature.
Provide ONLY the new jira title
Do NOT provide any other output.

Original Feature Request: {{.FeatureRequest}}

User Story: 
{{.UserStory}}
```

## Commands

### Generate Template Files

Generate example template files in a directory:

```bash
mcq templates generate /path/to/templates
```

### Validate Templates

Validate that all templates are syntactically correct:

```bash
mcq templates validate
```

### List Available Prompt Types

List all available prompt types and their template file names:

```bash
mcq templates list
```

## Template Syntax

Templates use Go's `html/template` syntax. Key features:

- `{{.VariableName}}` - Output a variable
- `{{if .Condition}}...{{end}}` - Conditional blocks
- `{{range .Items}}...{{end}}` - Loop over slices
- `{{/* comment */}}` - Comments
- `{{template "name" .}}` - Include other templates

## Error Handling

If template loading or execution fails, the system will:

1. Log the error
2. Fall back to the built-in default prompts
3. Continue operation normally

This ensures the application remains functional even with malformed templates.

## Development

For development, you can:

1. Set `MCQ_PROMPTS_DIR` to your development directory
2. Modify template files
3. Run `mcq templates validate` to check syntax
4. Test with actual commands

Templates are reloaded each time the application starts, so changes take effect immediately.
