package ai

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stevemcquaid/mcq/pkg/logger"
)

// TemplateData holds all available variables for prompt templates
type TemplateData struct {
	// Core request data
	FeatureRequest      string
	UserStory           string
	OriginalDescription string

	// Repository context
	RepositoryContext *RepoContext

	// Convenience fields for direct access
	ProjectName        string
	ModulePath         string
	GoVersion          string
	ProjectType        string
	Readme             string
	RecentCommits      []string
	Dependencies       []string
	DirectoryStructure string
	ConfigFiles        map[string]string

	// Template utilities
	Now time.Time
}

// TemplateManager handles loading and caching of prompt templates
type TemplateManager struct {
	templates  map[PromptType]*template.Template
	promptsDir string
}

var templateManager *TemplateManager

// GetTemplateManager returns the singleton template manager
func GetTemplateManager() *TemplateManager {
	if templateManager == nil {
		templateManager = &TemplateManager{
			templates:  make(map[PromptType]*template.Template),
			promptsDir: os.Getenv("MCQ_PROMPTS_DIR"),
		}
	}
	return templateManager
}

// LoadTemplates loads all prompt templates from the configured directory
func (tm *TemplateManager) LoadTemplates() error {
	// Clear existing templates
	tm.templates = make(map[PromptType]*template.Template)

	// If no custom directory is set, use default embedded templates
	if tm.promptsDir == "" {
		return tm.loadDefaultTemplates()
	}

	// Load templates from custom directory
	return tm.loadCustomTemplates()
}

// loadDefaultTemplates loads the default embedded templates
func (tm *TemplateManager) loadDefaultTemplates() error {
	// Create default templates for each prompt type
	for _, promptType := range []PromptType{PromptTypeUserStory, PromptTypeTitleExtraction, PromptTypeDescriptionImprovement, PromptTypeDescriptionFromTitle} {
		tmpl, err := tm.createDefaultTemplate(promptType)
		if err != nil {
			return fmt.Errorf("failed to create default template for %s: %w", promptType, err)
		}
		tm.templates[promptType] = tmpl
	}

	logger.LogBasic("Loaded default prompt templates")
	return nil
}

// loadCustomTemplates loads templates from the custom directory
func (tm *TemplateManager) loadCustomTemplates() error {
	if _, err := os.Stat(tm.promptsDir); os.IsNotExist(err) {
		return fmt.Errorf("prompts directory does not exist: %s", tm.promptsDir)
	}

	// Load each prompt type template
	for _, promptType := range []PromptType{PromptTypeUserStory, PromptTypeTitleExtraction, PromptTypeDescriptionImprovement, PromptTypeDescriptionFromTitle} {
		templateFile := filepath.Join(tm.promptsDir, string(promptType)+".tpl")

		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			logger.LogBasic("Template file not found, using default", "file", templateFile)
			// Fall back to default template
			tmpl, err := tm.createDefaultTemplate(promptType)
			if err != nil {
				return fmt.Errorf("failed to create default template for %s: %w", promptType, err)
			}
			tm.templates[promptType] = tmpl
			continue
		}

		// Create template with custom function map
		tmpl := template.New(filepath.Base(templateFile)).Funcs(template.FuncMap{
			"formatContext": formatContextForTemplate,
		})

		tmpl, err := tmpl.ParseFiles(templateFile)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", templateFile, err)
		}

		tm.templates[promptType] = tmpl
		logger.LogBasic("Loaded custom template", "file", templateFile)
	}

	logger.LogBasic("Loaded custom prompt templates", "directory", tm.promptsDir)
	return nil
}

// createDefaultTemplate creates a default template for the given prompt type
func (tm *TemplateManager) createDefaultTemplate(promptType PromptType) (*template.Template, error) {
	var templateContent string

	switch promptType {
	case PromptTypeUserStory:
		templateContent = `Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations. Provide ONLY the user story. 


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
{{end}}`

	case PromptTypeTitleExtraction:
		templateContent = `Create a NEW concise, clear title (maximum 100 characters) for a Jira issue from the following user story and old title. The new title should be action-oriented and summarize the main goal or feature.
Provide ONLY the new jira title
Do NOT provide any other output.

Original Feature Request: {{.FeatureRequest}}

User Story: 
{{.UserStory}}`

	case PromptTypeDescriptionImprovement:
		templateContent = `Improve the following Jira issue description. Make it:
1. More comprehensive and detailed
2. Better structured and readable
3. Include proper user story format if missing
4. Add acceptance criteria if not present
5. Add technical considerations
6. Ensure it follows best practices for user stories

Preserve the existing intent and structure, but enhance clarity, completeness, and professionalism.

Original Description:
{{.OriginalDescription}}
{{if .RepositoryContext}}
{{formatContext .RepositoryContext}}
{{end}}`

	case PromptTypeDescriptionFromTitle:
		templateContent = `Create a comprehensive user story description from the following Jira issue title.

The description should:
1. Follow the user story format: "As a [user type], I want [goal] so that [benefit]"
2. Be detailed and specific, not just repeating the title
3. Include acceptance criteria
4. Include technical considerations
5. Be comprehensive and well-structured

Title: {{.OriginalDescription}}
{{if .RepositoryContext}}
{{formatContext .RepositoryContext}}
{{end}}`
	}

	// Create template with repository context helper
	tmpl := template.New(string(promptType))
	tmpl = tmpl.Funcs(template.FuncMap{
		"formatContext": formatContextForTemplate,
	})

	return tmpl.Parse(templateContent)
}

// GeneratePromptFromTemplate generates a prompt using the template system
func (tm *TemplateManager) GeneratePromptFromTemplate(promptType PromptType, data TemplateData) (string, error) {
	tmpl, exists := tm.templates[promptType]

	if !exists {
		return "", fmt.Errorf("template not found for prompt type: %s", promptType)
	}

	// Prepare template data
	data = prepareTemplateData(data)

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// prepareTemplateData populates convenience fields from repository context
func prepareTemplateData(data TemplateData) TemplateData {
	data.Now = time.Now()
	if data.RepositoryContext != nil {
		ctx := data.RepositoryContext
		data.ProjectName = ctx.ProjectName
		data.ModulePath = ctx.ModulePath
		data.GoVersion = ctx.GoVersion
		data.ProjectType = ctx.ProjectType
		data.Readme = ctx.Readme
		data.RecentCommits = ctx.RecentCommits
		data.Dependencies = ctx.Dependencies
		data.DirectoryStructure = ctx.DirectoryStructure
		data.ConfigFiles = ctx.ConfigFiles
	}
	return data
}

// ValidateTemplates validates all loaded templates
func (tm *TemplateManager) ValidateTemplates() error {
	for promptType, tmpl := range tm.templates {
		// Test template with sample data
		testData := TemplateData{
			FeatureRequest: "Test feature request",
			UserStory:      "Test user story",
			Now:            time.Now(),
		}

		var result strings.Builder
		if err := tmpl.Execute(&result, testData); err != nil {
			return fmt.Errorf("template validation failed for %s: %w", promptType, err)
		}
	}

	return nil
}

// GetTemplateFile returns the expected template file path for a prompt type
func (tm *TemplateManager) GetTemplateFile(promptType PromptType) string {
	if tm.promptsDir == "" {
		return ""
	}
	return filepath.Join(tm.promptsDir, string(promptType)+".tpl")
}

// ReloadTemplates reloads all templates (useful for development)
func (tm *TemplateManager) ReloadTemplates() error {
	logger.LogBasic("Reloading prompt templates")
	return tm.LoadTemplates()
}

// formatContextForTemplate formats repository context for use in templates
func formatContextForTemplate(ctx *RepoContext) string {
	if ctx == nil {
		return ""
	}
	return formatContextForPrompt(ctx)
}
