package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stevemcquaid/mcq/pkg/ai"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage AI prompt templates",
	Long:  `Manage AI prompt templates for customizing AI behavior`,
}

var templatesValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate prompt templates",
	Long:  `Validate that all prompt templates are syntactically correct and can be executed`,
	Run: func(cmd *cobra.Command, args []string) {
		tm := ai.GetTemplateManager()

		// Load templates
		if err := tm.LoadTemplates(); err != nil {
			fmt.Printf("❌ Failed to load templates: %v\n", err)
			os.Exit(1)
		}

		// Validate templates
		if err := tm.ValidateTemplates(); err != nil {
			fmt.Printf("❌ Template validation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ All templates are valid")
	},
}

var templatesGenerateCmd = &cobra.Command{
	Use:   "generate [directory]",
	Short: "Generate example template files",
	Long:  `Generate example template files in the specified directory (or current directory if not specified)`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputDir := "."
		if len(args) > 0 {
			outputDir = args[0]
		}

		// Create directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			fmt.Printf("❌ Failed to create directory %s: %v\n", outputDir, err)
			os.Exit(1)
		}

		// Generate template files for each prompt type
		for _, promptType := range []ai.PromptType{ai.PromptTypeUserStory, ai.PromptTypeTitleExtraction} {
			templateFile := filepath.Join(outputDir, string(promptType)+".tpl")

			if _, err := os.Stat(templateFile); err == nil {
				fmt.Printf("⚠️  Template file already exists: %s\n", templateFile)
				continue
			}

			if err := generateTemplateFile(templateFile, promptType); err != nil {
				fmt.Printf("❌ Failed to generate template %s: %v\n", templateFile, err)
				os.Exit(1)
			}

			fmt.Printf("✅ Generated template: %s\n", templateFile)
		}

		fmt.Printf("\n📁 Template files generated in: %s\n", outputDir)
		fmt.Println("💡 Set MCQ_PROMPTS_DIR environment variable to use these templates:")
		fmt.Printf("   export MCQ_PROMPTS_DIR=%s\n", outputDir)
	},
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available prompt types",
	Long:  `List all available prompt types and their template file names`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available prompt types:")
		fmt.Println()

		for _, promptType := range []ai.PromptType{ai.PromptTypeUserStory, ai.PromptTypeTitleExtraction} {
			fmt.Printf("• %s\n", promptType)
			fmt.Printf("  Template file: %s.tpl\n", promptType)
			fmt.Printf("  Description: %s\n", getPromptTypeDescription(promptType))
			fmt.Println()
		}
	},
}

func init() {
	RootCmd.AddCommand(templatesCmd)
	templatesCmd.AddCommand(templatesValidateCmd)
	templatesCmd.AddCommand(templatesGenerateCmd)
	templatesCmd.AddCommand(templatesListCmd)
}

// generateTemplateFile generates a template file for the given prompt type
func generateTemplateFile(filePath string, promptType ai.PromptType) error {
	var content string

	switch promptType {
	case ai.PromptTypeUserStory:
		content = `{{/* 
User Story Generation Template
Available variables:
- .FeatureRequest: The user's feature request
- .RepositoryContext: Repository information (if available)
- .ProjectName: Project name from go.mod
- .ModulePath: Module path from go.mod
- .GoVersion: Go version from go.mod
- .ProjectType: Detected project type
- .Readme: README content
- .RecentCommits: Recent commit messages
- .Dependencies: Go dependencies
- .DirectoryStructure: Directory structure
- .ConfigFiles: Configuration files content
- .Now: Current timestamp
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
{{end}}`

	case ai.PromptTypeTitleExtraction:
		content = `{{/* 
Title Extraction Template
Available variables:
- .FeatureRequest: The original feature request
- .UserStory: The generated user story
- .Now: Current timestamp
*/}}
Create a NEW concise, clear title (maximum 100 characters) for a Jira issue from the following user story and old title. The new title should be action-oriented and summarize the main goal or feature.
Provide ONLY the new jira title
Do NOT provide any other output.

Original Feature Request: {{.FeatureRequest}}

User Story: 
{{.UserStory}}`
	}

	return os.WriteFile(filePath, []byte(content), 0o644)
}

// getPromptTypeDescription returns a description for the prompt type
func getPromptTypeDescription(promptType ai.PromptType) string {
	switch promptType {
	case ai.PromptTypeUserStory:
		return "Generates detailed user stories from feature requests"
	case ai.PromptTypeTitleExtraction:
		return "Extracts concise titles from user stories for JIRA issues"
	default:
		return "Unknown prompt type"
	}
}
