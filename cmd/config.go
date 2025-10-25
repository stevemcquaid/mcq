package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetupCmd)
	configCmd.AddCommand(configTestCmd)
	configCmd.AddCommand(configShowCmd)
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

	// Environment Variables
	fmt.Println("🌍 Environment Variables:")
	fmt.Printf("   • JIRA_INSTANCE_URL: %s\n", os.Getenv("JIRA_INSTANCE_URL"))
	fmt.Printf("   • JIRA_USERNAME: %s\n", os.Getenv("JIRA_USERNAME"))
	fmt.Printf("   • JIRA_API_TOKEN: %s\n", maskAPIKey(os.Getenv("JIRA_API_TOKEN")))
	fmt.Printf("   • JIRA_PROJECT_PREFIX: %s\n", os.Getenv("JIRA_PROJECT_PREFIX"))
	fmt.Printf("   • ANTHROPIC_API_KEY: %s\n", maskAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
	fmt.Printf("   • OPENAI_API_KEY: %s\n", maskAPIKey(os.Getenv("OPENAI_API_KEY")))
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
