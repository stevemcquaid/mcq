package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help information",
	Long:  `Display comprehensive help and examples for the MCQ CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		showHelp()
	},
}

// examplesCmd represents the examples command
var examplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show usage examples",
	Long:  `Display practical examples of how to use the MCQ CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		showExamples()
	},
}

// workflowCmd represents the workflow command
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Show recommended workflows",
	Long:  `Display recommended workflows and best practices for using the MCQ CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		showWorkflows()
	},
}

func init() {
	RootCmd.AddCommand(helpCmd)
	RootCmd.AddCommand(examplesCmd)
	RootCmd.AddCommand(workflowCmd)
}

// showHelp displays comprehensive help information
func showHelp() {
	fmt.Println("🔧 MCQ - Makefile Helper & Development Tools")
	fmt.Println("=============================================")
	fmt.Println()
	fmt.Println("MCQ is a powerful CLI tool that helps streamline your development workflow")
	fmt.Println("with JIRA integration and AI-powered user story generation.")
	fmt.Println()

	// Available Commands
	fmt.Println("📋 Available Commands:")
	fmt.Println("----------------------")
	fmt.Println()
	fmt.Println("🔧 Configuration:")
	fmt.Println("  config setup     Interactive configuration setup")
	fmt.Println("  config test      Test current configuration")
	fmt.Println("  config show      Show current configuration")
	fmt.Println()
	fmt.Println("📋 JIRA Commands:")
	fmt.Println("  jira show <key>  Display detailed JIRA issue information")
	fmt.Println("  jira new <story> Create new JIRA issue from user story")
	fmt.Println()
	fmt.Println("🤖 AI Commands:")
	fmt.Println("  ai jira <story>  Convert feature request to user story")
	fmt.Println()
	fmt.Println("❓ Help Commands:")
	fmt.Println("  help             Show this help information")
	fmt.Println("  examples         Show usage examples")
	fmt.Println("  workflow         Show recommended workflows")
	fmt.Println()

	// Quick Start
	fmt.Println("🚀 Quick Start:")
	fmt.Println("---------------")
	fmt.Println("1. Run 'mcq config setup' to configure JIRA and AI settings")
	fmt.Println("2. Run 'mcq config test' to verify your configuration")
	fmt.Println("3. Try 'mcq ai jira \"Add dark mode\"' to generate a user story")
	fmt.Println("4. Use 'mcq jira new \"Add dark mode\"' to create a JIRA issue")
	fmt.Println()

	// Environment Variables
	fmt.Println("🌍 Environment Variables:")
	fmt.Println("-------------------------")
	fmt.Println("JIRA_INSTANCE_URL    Your JIRA instance URL")
	fmt.Println("JIRA_USERNAME        Your JIRA username/email")
	fmt.Println("JIRA_API_TOKEN       Your JIRA API token")
	fmt.Println("JIRA_PROJECT_PREFIX  Your JIRA project prefix (e.g., PROJ)")
	fmt.Println("ANTHROPIC_API_KEY    Your Anthropic API key for Claude")
	fmt.Println("OPENAI_API_KEY       Your OpenAI API key for GPT models")
	fmt.Println()

	// More Information
	fmt.Println("📚 More Information:")
	fmt.Println("-------------------")
	fmt.Println("• Run 'mcq examples' for detailed usage examples")
	fmt.Println("• Run 'mcq workflow' for recommended workflows")
	fmt.Println("• Run 'mcq <command> --help' for command-specific help")
	fmt.Println()
}

// showExamples displays practical usage examples
func showExamples() {
	fmt.Println("📚 MCQ Usage Examples")
	fmt.Println("=====================")
	fmt.Println()

	// Configuration Examples
	fmt.Println("🔧 Configuration Examples:")
	fmt.Println("-------------------------")
	fmt.Println()
	fmt.Println("# Interactive setup")
	fmt.Println("mcq config setup")
	fmt.Println()
	fmt.Println("# Test configuration")
	fmt.Println("mcq config test")
	fmt.Println()
	fmt.Println("# Show current configuration")
	fmt.Println("mcq config show")
	fmt.Println()

	// JIRA Examples
	fmt.Println("📋 JIRA Examples:")
	fmt.Println("-----------------")
	fmt.Println()
	fmt.Println("# Show issue details")
	fmt.Println("mcq jira show PROJ-123")
	fmt.Println("mcq jira show 123  # Uses JIRA_PROJECT_PREFIX")
	fmt.Println()
	fmt.Println("# Create new issue from user story")
	fmt.Println("mcq jira new \"Add dark mode to the application\"")
	fmt.Println("mcq jira new -- Add dark mode to the application")
	fmt.Println()
	fmt.Println("# Create issue with specific AI model")
	fmt.Println("mcq jira new --model claude \"Improve user login process\"")
	fmt.Println()
	fmt.Println("# Create issue with context gathering")
	fmt.Println("mcq jira new --auto-context \"Add user authentication\"")
	fmt.Println()

	// AI Examples
	fmt.Println("🤖 AI Examples:")
	fmt.Println("---------------")
	fmt.Println()
	fmt.Println("# Generate user story")
	fmt.Println("mcq ai jira \"Add dark mode to the application\"")
	fmt.Println()
	fmt.Println("# Generate with specific model")
	fmt.Println("mcq ai jira --model gpt-5 \"Improve user login process\"")
	fmt.Println()
	fmt.Println("# Generate with context")
	fmt.Println("mcq ai jira --include-readme --include-commits \"Add user authentication\"")
	fmt.Println()
	fmt.Println("# Generate without context")
	fmt.Println("mcq ai jira --no-context \"Add user authentication\"")
	fmt.Println()

	// Advanced Examples
	fmt.Println("⚡ Advanced Examples:")
	fmt.Println("--------------------")
	fmt.Println()
	fmt.Println("# Verbose output for debugging")
	fmt.Println("mcq jira new --verbosity 3 \"Add dark mode\"")
	fmt.Println()
	fmt.Println("# Custom context configuration")
	fmt.Println("mcq ai jira --include-readme --include-go-mod --max-commits 5 \"Add feature\"")
	fmt.Println()
	fmt.Println("# Batch processing (future feature)")
	fmt.Println("mcq batch jira features.txt")
	fmt.Println()

	// Troubleshooting Examples
	fmt.Println("🔧 Troubleshooting Examples:")
	fmt.Println("----------------------------")
	fmt.Println()
	fmt.Println("# Test configuration")
	fmt.Println("mcq config test")
	fmt.Println()
	fmt.Println("# Show detailed error information")
	fmt.Println("mcq jira show PROJ-123 --verbose")
	fmt.Println()
	fmt.Println("# Check available AI models")
	fmt.Println("mcq ai models")
	fmt.Println()
}

// showWorkflows displays recommended workflows
func showWorkflows() {
	fmt.Println("🔄 Recommended Workflows")
	fmt.Println("========================")
	fmt.Println()

	// Initial Setup Workflow
	fmt.Println("🚀 Initial Setup Workflow:")
	fmt.Println("--------------------------")
	fmt.Println("1. Clone your repository")
	fmt.Println("2. Run 'mcq config setup' to configure JIRA and AI")
	fmt.Println("3. Run 'mcq config test' to verify everything works")
	fmt.Println("4. Run 'mcq examples' to see what's possible")
	fmt.Println()

	// Daily Development Workflow
	fmt.Println("💻 Daily Development Workflow:")
	fmt.Println("-----------------------------")
	fmt.Println("1. Start with a vague idea: 'Add user authentication'")
	fmt.Println("2. Generate user story: 'mcq ai jira \"Add user authentication\"'")
	fmt.Println("3. Review and refine the generated user story")
	fmt.Println("4. Create JIRA issue: 'mcq jira new \"Add user authentication\"'")
	fmt.Println("5. Work on the issue and update JIRA as needed")
	fmt.Println()

	// Feature Development Workflow
	fmt.Println("🎯 Feature Development Workflow:")
	fmt.Println("--------------------------------")
	fmt.Println("1. Generate user story with context:")
	fmt.Println("   'mcq ai jira --auto-context \"Add dark mode\"'")
	fmt.Println("2. Create JIRA issue with AI-generated title:")
	fmt.Println("   'mcq jira new --auto-context \"Add dark mode\"'")
	fmt.Println("3. Review the generated content and make adjustments")
	fmt.Println("4. Assign the issue to yourself or team members")
	fmt.Println("5. Start development and update JIRA with progress")
	fmt.Println()

	// Team Collaboration Workflow
	fmt.Println("👥 Team Collaboration Workflow:")
	fmt.Println("------------------------------")
	fmt.Println("1. Product manager creates user stories using AI")
	fmt.Println("2. Stories are automatically created in JIRA")
	fmt.Println("3. Development team reviews and refines stories")
	fmt.Println("4. Stories are assigned to developers")
	fmt.Println("5. Progress is tracked in JIRA")
	fmt.Println()

	// Context Management Workflow
	fmt.Println("📁 Context Management Workflow:")
	fmt.Println("------------------------------")
	fmt.Println("1. Save context profiles for different project types:")
	fmt.Println("   'mcq context save go-project'")
	fmt.Println("2. Load context profiles when switching projects:")
	fmt.Println("   'mcq context load go-project'")
	fmt.Println("3. Use auto-context for new projects:")
	fmt.Println("   'mcq ai jira --auto-context \"New feature\"'")
	fmt.Println()

	// Troubleshooting Workflow
	fmt.Println("🔧 Troubleshooting Workflow:")
	fmt.Println("---------------------------")
	fmt.Println("1. Check configuration: 'mcq config show'")
	fmt.Println("2. Test configuration: 'mcq config test'")
	fmt.Println("3. Check specific command: 'mcq <command> --help'")
	fmt.Println("4. Enable verbose output: 'mcq <command> --verbose'")
	fmt.Println("5. Check examples: 'mcq examples'")
	fmt.Println()

	// Best Practices
	fmt.Println("✨ Best Practices:")
	fmt.Println("-----------------")
	fmt.Println("• Always test your configuration after setup")
	fmt.Println("• Use context gathering for better AI results")
	fmt.Println("• Review AI-generated content before creating issues")
	fmt.Println("• Keep your API keys secure and rotate them regularly")
	fmt.Println("• Use descriptive feature requests for better results")
	fmt.Println("• Regularly update your context profiles")
	fmt.Println()
}

// showCommandHelp displays help for a specific command
func showCommandHelp(command string) {
	switch command {
	case "jira":
		showJiraHelp()
	case "ai":
		showAIHelp()
	case "config":
		showConfigHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'mcq help' to see available commands")
	}
}

// showJiraHelp displays JIRA-specific help
func showJiraHelp() {
	fmt.Println("📋 JIRA Command Help")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println("The JIRA commands help you interact with your JIRA instance")
	fmt.Println("for issue management and creation.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  show <issue-key>  Display detailed issue information")
	fmt.Println("  new <story>       Create new issue from user story")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mcq jira show PROJ-123")
	fmt.Println("  mcq jira new \"Add dark mode\"")
	fmt.Println()
}

// showAIHelp displays AI-specific help
func showAIHelp() {
	fmt.Println("🤖 AI Command Help")
	fmt.Println("==================")
	fmt.Println()
	fmt.Println("The AI commands help you generate user stories and other")
	fmt.Println("content using artificial intelligence.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  jira <story>      Convert feature request to user story")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mcq ai jira \"Add dark mode\"")
	fmt.Println("  mcq ai jira --model claude \"Improve login\"")
	fmt.Println()
}

// showConfigHelp displays configuration-specific help
func showConfigHelp() {
	fmt.Println("🔧 Configuration Command Help")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("The configuration commands help you set up and manage")
	fmt.Println("your JIRA and AI settings.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  setup             Interactive configuration setup")
	fmt.Println("  test              Test current configuration")
	fmt.Println("  show              Show current configuration")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mcq config setup")
	fmt.Println("  mcq config test")
	fmt.Println("  mcq config show")
	fmt.Println()
}
