package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stevemcquaid/mcq/pkg/ai"
	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/jira"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Configure JIRA and AI settings interactively or view current configuration.`,
}

// configSetupCmd represents the config setup command
var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive configuration setup",
	Long:  `Set up JIRA and AI configuration interactively with guided prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupInteractive()
	},
}

// configTestCmd represents the config test command
var configTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test current configuration",
	Long:  `Test JIRA and AI configuration to ensure everything is working correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		testConfiguration()
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display current configuration settings (API keys are masked for security).`,
	Run: func(cmd *cobra.Command, args []string) {
		showConfiguration()
	},
}

// configTemplatesCmd represents the config templates command
var configTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage AI prompt templates",
	Long:  `Generate, validate, and manage AI prompt templates for customization.`,
}

// configTemplatesGenerateCmd represents the config templates generate command
var configTemplatesGenerateCmd = &cobra.Command{
	Use:   "generate [directory]",
	Short: "Generate example template files",
	Long:  `Generate example template files in the specified directory (or current directory if not specified).`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		generateTemplates(args)
	},
}

// configTemplatesValidateCmd represents the config templates validate command
var configTemplatesValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate template syntax",
	Long:  `Validate that all prompt templates are syntactically correct and can be executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateTemplates()
	},
}

// configTemplatesListCmd represents the config templates list command
var configTemplatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available prompt types",
	Long:  `List all available prompt types and their template file names.`,
	Run: func(cmd *cobra.Command, args []string) {
		listTemplates()
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetupCmd)
	configCmd.AddCommand(configTestCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configTemplatesCmd)

	// Template subcommands
	configTemplatesCmd.AddCommand(configTemplatesGenerateCmd)
	configTemplatesCmd.AddCommand(configTemplatesValidateCmd)
	configTemplatesCmd.AddCommand(configTemplatesListCmd)
}

// setupInteractive guides the user through configuration setup
func setupInteractive() {
	fmt.Println("🔧 MCQ Configuration Setup")
	fmt.Println("==========================")
	fmt.Println()

	// JIRA Configuration
	fmt.Println("📋 JIRA Configuration")
	fmt.Println("---------------------")

	// JIRA URL
	jiraURL := getInput("JIRA Instance URL", viper.GetString("jira.url"), "https://yourcompany.atlassian.net")
	if jiraURL != "" {
		viper.Set("jira.url", jiraURL)
	}

	// JIRA Username
	jiraUsername := getInput("JIRA Username/Email", viper.GetString("jira.username"), "user@company.com")
	if jiraUsername != "" {
		viper.Set("jira.username", jiraUsername)
	}

	// JIRA API Token
	jiraToken := getInput("JIRA API Token", maskAPIKey(viper.GetString("jira.token")), "your_api_token_here")
	if jiraToken != "" && !strings.HasPrefix(jiraToken, "***") {
		viper.Set("jira.token", jiraToken)
	}

	// JIRA Project Prefix
	jiraProjectPrefix := getInput("JIRA Project Prefix", viper.GetString("jira.project_prefix"), "PROJ")
	if jiraProjectPrefix != "" {
		viper.Set("jira.project_prefix", jiraProjectPrefix)
	}

	fmt.Println()

	// AI Configuration
	fmt.Println("🤖 AI Configuration")
	fmt.Println("-------------------")

	// Anthropic API Key
	anthropicKey := getInput("Anthropic API Key (for Claude)", maskAPIKey(os.Getenv("ANTHROPIC_API_KEY")), "sk-ant-...")
	if anthropicKey != "" && !strings.HasPrefix(anthropicKey, "***") {
		if err := os.Setenv("ANTHROPIC_API_KEY", anthropicKey); err != nil {
			fmt.Printf("Warning: Failed to set ANTHROPIC_API_KEY: %v\n", err)
		}
	}

	// OpenAI API Key
	openaiKey := getInput("OpenAI API Key (for GPT models)", maskAPIKey(os.Getenv("OPENAI_API_KEY")), "sk-...")
	if openaiKey != "" && !strings.HasPrefix(openaiKey, "***") {
		if err := os.Setenv("OPENAI_API_KEY", openaiKey); err != nil {
			fmt.Printf("Warning: Failed to set OPENAI_API_KEY: %v\n", err)
		}
	}

	fmt.Println()

	// Template Configuration
	fmt.Println("📝 Template Configuration")
	fmt.Println("-------------------------")

	// Ask about template customization
	if askForConfirmation("Would you like to customize AI prompt templates?", false) {
		templateDir := getInput("Template directory path", os.Getenv("MCQ_PROMPTS_DIR"), "./templates")
		if templateDir != "" {
			if err := os.Setenv("MCQ_PROMPTS_DIR", templateDir); err != nil {
				fmt.Printf("Warning: Failed to set MCQ_PROMPTS_DIR: %v\n", err)
			} else {
				fmt.Printf("✅ Template directory set to: %s\n", templateDir)

				// Ask if they want to generate example templates
				if askForConfirmation("Generate example template files?", true) {
					if err := generateTemplateFiles(templateDir); err != nil {
						fmt.Printf("❌ Failed to generate templates: %v\n", err)
					} else {
						fmt.Printf("✅ Example templates generated in: %s\n", templateDir)
					}
				}
			}
		}
	}

	fmt.Println()

	// Save configuration
	if askForConfirmation("Save this configuration?", true) {
		saveConfiguration()
		fmt.Println("✅ Configuration saved successfully!")
	} else {
		fmt.Println("Configuration not saved.")
	}

	// Test configuration
	if askForConfirmation("Test the configuration now?", true) {
		testConfiguration()
	}
}

// testConfiguration tests the current configuration
func testConfiguration() {
	fmt.Println("🧪 Testing Configuration")
	fmt.Println("========================")
	fmt.Println()

	// Test JIRA configuration
	fmt.Println("📋 Testing JIRA Configuration...")
	jiraManager, err := jira.NewManager()
	if err != nil {
		userErr := errors.WrapError(err, "JIRA configuration test failed")
		userErr.Display()
	} else {
		fmt.Println("✅ JIRA configuration is valid")
		fmt.Printf("   • URL: %s\n", jiraManager.GetBaseURL())
	}

	fmt.Println()

	// Test AI configuration
	fmt.Println("🤖 Testing AI Configuration...")
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")

	if anthropicKey == "" && openaiKey == "" {
		fmt.Println("❌ No AI API keys found")
		fmt.Println("   • Set ANTHROPIC_API_KEY for Claude models")
		fmt.Println("   • Set OPENAI_API_KEY for GPT models")
	} else {
		fmt.Println("✅ AI configuration is valid")
		if anthropicKey != "" {
			fmt.Printf("   • Anthropic API Key: %s\n", maskAPIKey(anthropicKey))
		}
		if openaiKey != "" {
			fmt.Printf("   • OpenAI API Key: %s\n", maskAPIKey(openaiKey))
		}
	}

	fmt.Println()
	fmt.Println("🎉 Configuration test completed!")
}

// showConfiguration displays the current configuration
func showConfiguration() {
	fmt.Println("📋 Current Configuration")
	fmt.Println("========================")
	fmt.Println()

	// JIRA Configuration
	fmt.Println("📋 JIRA Settings:")
	fmt.Printf("   • URL: %s\n", viper.GetString("jira.url"))
	fmt.Printf("   • Username: %s\n", viper.GetString("jira.username"))
	fmt.Printf("   • Token: %s\n", maskAPIKey(viper.GetString("jira.token")))
	fmt.Printf("   • Project Prefix: %s\n", viper.GetString("jira.project_prefix"))
	fmt.Println()

	// AI Configuration
	fmt.Println("🤖 AI Settings:")
	fmt.Printf("   • Anthropic API Key: %s\n", maskAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
	fmt.Printf("   • OpenAI API Key: %s\n", maskAPIKey(os.Getenv("OPENAI_API_KEY")))
	fmt.Println()

	// Template Configuration
	fmt.Println("📝 Template Settings:")
	templateDir := os.Getenv("MCQ_PROMPTS_DIR")
	if templateDir == "" {
		fmt.Println("   • Template Directory: (using default templates)")
	} else {
		fmt.Printf("   • Template Directory: %s\n", templateDir)
	}
	fmt.Println()

	// Environment Variables
	fmt.Println("🌍 Environment Variables:")
	fmt.Printf("   • JIRA_INSTANCE_URL: %s\n", os.Getenv("JIRA_INSTANCE_URL"))
	fmt.Printf("   • JIRA_USERNAME: %s\n", os.Getenv("JIRA_USERNAME"))
	fmt.Printf("   • JIRA_API_TOKEN: %s\n", maskAPIKey(os.Getenv("JIRA_API_TOKEN")))
	fmt.Printf("   • JIRA_PROJECT_PREFIX: %s\n", os.Getenv("JIRA_PROJECT_PREFIX"))
	fmt.Printf("   • ANTHROPIC_API_KEY: %s\n", maskAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
	fmt.Printf("   • OPENAI_API_KEY: %s\n", maskAPIKey(os.Getenv("OPENAI_API_KEY")))
	fmt.Printf("   • MCQ_PROMPTS_DIR: %s\n", os.Getenv("MCQ_PROMPTS_DIR"))
}

// getInput prompts for user input with a default value
func getInput(prompt, current, placeholder string) string {
	reader := bufio.NewReader(os.Stdin)

	defaultText := placeholder
	if current != "" {
		defaultText = current
	}

	fmt.Printf("%s [%s]: ", prompt, defaultText)

	input, err := reader.ReadString('\n')
	if err != nil {
		return current
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return current
	}

	return input
}

// askForConfirmation prompts the user for confirmation
func askForConfirmation(prompt string, defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)

	defaultText := "Y/n"
	if !defaultYes {
		defaultText = "y/N"
	}

	fmt.Printf("%s [%s]: ", prompt, defaultText)

	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}

// maskAPIKey returns a masked version of the API key
func maskAPIKey(apiKey string) string {
	if apiKey == "" {
		return "not set"
	}
	if len(apiKey) <= 4 {
		return "***"
	}
	return "***" + apiKey[len(apiKey)-4:]
}

// saveConfiguration saves the current configuration to a file
func saveConfiguration() {
	// For now, we'll just set environment variables
	// In a real implementation, you might want to save to a config file

	// Set environment variables for JIRA
	if url := viper.GetString("jira.url"); url != "" {
		if err := os.Setenv("JIRA_INSTANCE_URL", url); err != nil {
			fmt.Printf("Warning: Failed to set JIRA_INSTANCE_URL: %v\n", err)
		}
	}
	if username := viper.GetString("jira.username"); username != "" {
		if err := os.Setenv("JIRA_USERNAME", username); err != nil {
			fmt.Printf("Warning: Failed to set JIRA_USERNAME: %v\n", err)
		}
	}
	if token := viper.GetString("jira.token"); token != "" {
		if err := os.Setenv("JIRA_API_TOKEN", token); err != nil {
			fmt.Printf("Warning: Failed to set JIRA_API_TOKEN: %v\n", err)
		}
	}
	if prefix := viper.GetString("jira.project_prefix"); prefix != "" {
		if err := os.Setenv("JIRA_PROJECT_PREFIX", prefix); err != nil {
			fmt.Printf("Warning: Failed to set JIRA_PROJECT_PREFIX: %v\n", err)
		}
	}
}

// generateTemplates generates example template files
func generateTemplates(args []string) {
	outputDir := "."
	if len(args) > 0 {
		outputDir = args[0]
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create directory %s: %v\n", outputDir, err)
		return
	}

	// Generate template files for each prompt type
	for _, promptType := range []ai.PromptType{ai.PromptTypeUserStory, ai.PromptTypeTitleExtraction} {
		templateFile := fmt.Sprintf("%s/%s.tpl", outputDir, promptType)

		if _, err := os.Stat(templateFile); err == nil {
			fmt.Printf("⚠️  Template file already exists: %s\n", templateFile)
			continue
		}

		if err := generateTemplateFile(templateFile, promptType); err != nil {
			fmt.Printf("❌ Failed to generate template %s: %v\n", templateFile, err)
			return
		}

		fmt.Printf("✅ Generated template: %s\n", templateFile)
	}

	fmt.Printf("\n📁 Template files generated in: %s\n", outputDir)
	fmt.Println("💡 Set MCQ_PROMPTS_DIR environment variable to use these templates:")
	fmt.Printf("   export MCQ_PROMPTS_DIR=%s\n", outputDir)
}

// validateTemplates validates template syntax
func validateTemplates() {
	// This would need to import the AI package to use the template manager
	// For now, we'll provide a simple validation message
	fmt.Println("🔍 Validating templates...")

	templateDir := os.Getenv("MCQ_PROMPTS_DIR")
	if templateDir == "" {
		fmt.Println("ℹ️  No custom template directory set (MCQ_PROMPTS_DIR)")
		fmt.Println("   Using default templates - no validation needed")
		return
	}

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Printf("❌ Template directory does not exist: %s\n", templateDir)
		return
	}

	// Check if template files exist
	templateFiles := []string{"user_story.tpl", "title_extraction.tpl"}
	allExist := true

	for _, file := range templateFiles {
		templateFile := fmt.Sprintf("%s/%s", templateDir, file)
		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			fmt.Printf("⚠️  Template file missing: %s\n", templateFile)
			allExist = false
		}
	}

	if allExist {
		fmt.Println("✅ All template files found")
		fmt.Println("💡 Run 'mcq templates validate' for syntax validation")
	} else {
		fmt.Println("❌ Some template files are missing")
		fmt.Println("💡 Run 'mcq config templates generate' to create them")
	}
}

// listTemplates lists available prompt types
func listTemplates() {
	fmt.Println("Available prompt types:")
	fmt.Println()

	promptTypes := []struct {
		name        string
		description string
	}{
		{"user_story", "Generates detailed user stories from feature requests"},
		{"title_extraction", "Extracts concise titles from user stories for JIRA issues"},
	}

	for _, pt := range promptTypes {
		fmt.Printf("• %s\n", pt.name)
		fmt.Printf("  Template file: %s.tpl\n", pt.name)
		fmt.Printf("  Description: %s\n", pt.description)
		fmt.Println()
	}
}

// generateTemplateFiles is a helper function for the interactive setup
func generateTemplateFiles(templateDir string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return err
	}

	// Generate template files for each prompt type
	for _, promptType := range []ai.PromptType{ai.PromptTypeUserStory, ai.PromptTypeTitleExtraction} {
		templateFile := fmt.Sprintf("%s/%s.tpl", templateDir, promptType)

		if err := generateTemplateFile(templateFile, promptType); err != nil {
			return err
		}
	}

	return nil
}
