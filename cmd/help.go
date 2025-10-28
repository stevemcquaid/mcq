package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help information",
	Long:  `Display comprehensive help and examples for the MCQ CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			showCommandHelp(args[0])
		} else {
			showHelp()
		}
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

// quickCmd represents the quick reference command
var quickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Show quick reference",
	Long:  `Display a quick reference of all available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		showQuickReference()
	},
}

// commandsCmd represents the commands list command
var commandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "List all commands",
	Long:  `Display a categorized list of all available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		listAllCommands()
	},
}

func init() {
	RootCmd.AddCommand(helpCmd)
	RootCmd.AddCommand(examplesCmd)
	RootCmd.AddCommand(workflowCmd)
	RootCmd.AddCommand(quickCmd)
	RootCmd.AddCommand(commandsCmd)
}

// showCommandHelp shows help for a specific command
func showCommandHelp(cmdName string) {
	// Try to find and show specific command help
	if cmd := RootCmd.Commands(); cmd != nil {
		for _, c := range cmd {
			if c.Name() == cmdName || inAliases(c, cmdName) {
				c.Help()
				return
			}
			// Check subcommands
			for _, subc := range c.Commands() {
				if subc.Name() == cmdName || inAliases(subc, cmdName) {
					subc.Help()
					return
				}
			}
		}
	}

	// If not found, show general help
	fmt.Printf("No command found: %s\n\n", cmdName)
	showHelp()
}

// inAliases checks if a name is in the command's aliases
func inAliases(cmd *cobra.Command, name string) bool {
	for _, alias := range cmd.Aliases {
		if alias == name {
			return true
		}
	}
	return false
}

// showHelp displays comprehensive help information
func showHelp() {
	fmt.Println("🔧 MCQ - Makefile Helper & Development Tools")
	fmt.Println("=============================================")
	fmt.Println()
	fmt.Println("MCQ is a powerful CLI tool that helps streamline your development workflow")
	fmt.Println("with JIRA integration, AI-powered user story generation, and customizable prompt templates.")
	fmt.Println()

	// Available Commands
	fmt.Println("📋 Available Commands:")
	fmt.Println("----------------------")
	fmt.Println()
	fmt.Println("🔧 Configuration:")
	fmt.Println("  config setup     Interactive configuration setup")
	fmt.Println("  config test      Test current configuration")
	fmt.Println("  config show      Show current configuration")
	fmt.Println("  context test     Test repository context gathering")
	fmt.Println()
	fmt.Println("💻 Development:")
	fmt.Println("  build            Build the binary")
	fmt.Println("  test             Run all tests")
	fmt.Println("  fmt              Format code")
	fmt.Println("  lint             Run linters")
	fmt.Println("  deps             Manage dependencies")
	fmt.Println("  run              Run main.go")
	fmt.Println("  clean            Prep for commit")
	fmt.Println("  ci               Run CI checks")
	fmt.Println("  all              Run everything")
	fmt.Println()
	fmt.Println("📋 JIRA Commands:")
	fmt.Println("  jira show <key>     Display detailed JIRA issue information")
	fmt.Println("  jira new <story>    Create JIRA issue from feature request")
	fmt.Println("  jira new --dry-run  Generate user story without creating ticket")
	fmt.Println()
	fmt.Println("📝 Template Commands:")
	fmt.Println("  templates generate [dir]  Generate example template files")
	fmt.Println("  templates validate       Validate template syntax")
	fmt.Println("  templates list           List available prompt types")
	fmt.Println()
	fmt.Println("🐳 Docker:")
	fmt.Println("  docker build       Build docker image")
	fmt.Println("  docker run         Run docker container")
	fmt.Println("  docker push        Push to registry")
	fmt.Println()
	fmt.Println("📂 Git:")
	fmt.Println("  log                Pretty git log")
	fmt.Println("  gitclean           Clean working directory")
	fmt.Println()
	fmt.Println("⚙️ Setup:")
	fmt.Println("  setup              Install dependencies")
	fmt.Println("  install            Install the binary")
	fmt.Println("  version            Show version")
	fmt.Println()
	fmt.Println("❓ Help Commands:")
	fmt.Println("  help             Show this help information")
	fmt.Println("  examples         Show usage examples")
	fmt.Println("  workflow         Show recommended workflows")
	fmt.Println("  quick            Show quick reference")
	fmt.Println("  commands         List all commands categorized")
	fmt.Println()

	// Quick Start
	fmt.Println("🚀 Quick Start:")
	fmt.Println("---------------")
	fmt.Println("1. Run 'mcq config setup' to configure JIRA and AI settings")
	fmt.Println("2. Run 'mcq config test' to verify your configuration")
	fmt.Println("3. Try 'mcq jira new --dry-run \"Add dark mode\"' to generate a user story")
	fmt.Println("4. Use 'mcq jira new \"Add dark mode\"' to create a JIRA issue")
	fmt.Println()
	fmt.Println("📝 Template Customization:")
	fmt.Println("---------------------------")
	fmt.Println("1. Run 'mcq templates generate ./my-templates' to create example templates")
	fmt.Println("2. Set 'export MCQ_PROMPTS_DIR=./my-templates' to use custom templates")
	fmt.Println("3. Edit the .tpl files to customize AI prompts")
	fmt.Println("4. Run 'mcq templates validate' to check template syntax")
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
	fmt.Println("MCQ_PROMPTS_DIR      Directory containing custom prompt templates")
	fmt.Println()

	// More Information
	fmt.Println("📚 More Information:")
	fmt.Println("-------------------")
	fmt.Println("• Run 'mcq help' for comprehensive help")
	fmt.Println("• Run 'mcq examples' for detailed usage examples")
	fmt.Println("• Run 'mcq workflow' for recommended workflows")
	fmt.Println("• Run 'mcq quick' for a quick reference")
	fmt.Println("• Run 'mcq commands' to list all commands")
	fmt.Println("• Run 'mcq help <command>' for command-specific help")
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
	fmt.Println("# Test context gathering")
	fmt.Println("mcq context test")
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
	fmt.Println("# Generate user story without creating ticket")
	fmt.Println("mcq jira new --dry-run \"Add dark mode to the application\"")
	fmt.Println()
	fmt.Println("# Generate with specific model")
	fmt.Println("mcq jira new --dry-run --model gpt-5 \"Improve user login process\"")
	fmt.Println()
	fmt.Println("# Generate with context")
	fmt.Println("mcq jira new --dry-run --include-readme --include-commits \"Add user authentication\"")
	fmt.Println()
	fmt.Println("# Generate without context")
	fmt.Println("mcq jira new --dry-run --no-context \"Add user authentication\"")
	fmt.Println()

	// Advanced Examples
	fmt.Println("⚡ Advanced Examples:")
	fmt.Println("--------------------")
	fmt.Println()
	fmt.Println("# Verbose output for debugging")
	fmt.Println("mcq jira new --verbosity 3 \"Add dark mode\"")
	fmt.Println()
	fmt.Println("# Custom context configuration")
	fmt.Println("mcq jira new --dry-run --include-readme --include-go-mod --max-commits 5 \"Add feature\"")
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
	fmt.Println("2. Generate user story (dry-run): 'mcq jira new --dry-run \"Add user authentication\"'")
	fmt.Println("3. Review and refine the generated user story")
	fmt.Println("4. Create JIRA issue: 'mcq jira new \"Add user authentication\"'")
	fmt.Println("5. Work on the issue and update JIRA as needed")
	fmt.Println()

	// Feature Development Workflow
	fmt.Println("🎯 Feature Development Workflow:")
	fmt.Println("--------------------------------")
	fmt.Println("1. Generate user story with context (dry-run):")
	fmt.Println("   'mcq jira new --dry-run --auto-context \"Add dark mode\"'")
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
	fmt.Println("   'mcq jira new --dry-run --auto-context \"New feature\"'")
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

// showQuickReference displays a quick reference of all commands
func showQuickReference() {
	fmt.Println("⚡ Quick Reference - MCQ Commands")
	fmt.Println("==================================")
	fmt.Println()

	// Configuration Commands
	fmt.Println("🔧 CONFIGURATION:")
	fmt.Println("  mcq config setup     Interactive setup for JIRA and AI")
	fmt.Println("  mcq config test     Test your configuration")
	fmt.Println("  mcq config show     Display current configuration")
	fmt.Println("  mcq context test    Test context gathering")
	fmt.Println()

	// Development Commands
	fmt.Println("💻 DEVELOPMENT:")
	fmt.Println("  mcq build           Build the binary")
	fmt.Println("  mcq test            Run all tests")
	fmt.Println("  mcq test unit       Run unit tests only")
	fmt.Println("  mcq fmt             Format code")
	fmt.Println("  mcq lint            Run linters")
	fmt.Println("  mcq lint --fix      Run linters and fix issues")
	fmt.Println("  mcq vet             Run go vet")
	fmt.Println("  mcq deps            Tidy and vendor dependencies")
	fmt.Println("  mcq cover           Generate test coverage report")
	fmt.Println("  mcq run             Run main.go")
	fmt.Println()

	// JIRA Commands
	fmt.Println("📋 JIRA:")
	fmt.Println("  mcq jira show <key>     View issue details")
	fmt.Println("  mcq jira new <story>     Create issue from feature request")
	fmt.Println()

	// AI Commands
	fmt.Println("🤖 JIRA + AI:")
	fmt.Println("  mcq jira new <story>         Create JIRA issue")
	fmt.Println("  mcq jira new --dry-run        Generate without creating")
	fmt.Println("  mcq jira new --model claude   Specify AI model")
	fmt.Println("  mcq jira new --auto-context   Include repo context")
	fmt.Println()

	// Docker Commands
	fmt.Println("🐳 DOCKER:")
	fmt.Println("  mcq docker build      Build docker image")
	fmt.Println("  mcq docker run        Run docker container")
	fmt.Println("  mcq docker push       Push to registry")
	fmt.Println()

	// Git Commands
	fmt.Println("📂 GIT:")
	fmt.Println("  mcq log             Pretty git log")
	fmt.Println("  mcq gitclean        Clean working directory")
	fmt.Println()

	// CI/CD Commands
	fmt.Println("🚀 CI/CD:")
	fmt.Println("  mcq ci              Run CI checks")
	fmt.Println("  mcq all             Run everything")
	fmt.Println("  mcq clean           Prep for commit")
	fmt.Println()

	// Help Commands
	fmt.Println("❓ HELP:")
	fmt.Println("  mcq help             Show this help")
	fmt.Println("  mcq examples         Show usage examples")
	fmt.Println("  mcq workflow         Show workflows")
	fmt.Println("  mcq quick            Show this quick reference")
	fmt.Println("  mcq commands         List all commands")
	fmt.Println("  mcq version          Show version")
	fmt.Println()

	// Tips
	fmt.Println("💡 TIPS:")
	fmt.Println("  • Use 'mcq <command> --help' for detailed help")
	fmt.Println("  • Set environment variables for JIRA and AI")
	fmt.Println("  • Use --verbosity flag for debugging")
	fmt.Println()

	// Environment Variables
	fmt.Println("🌍 KEY ENV VARS:")
	fmt.Println("  JIRA_INSTANCE_URL      Your JIRA URL")
	fmt.Println("  JIRA_API_TOKEN         Your JIRA token")
	fmt.Println("  JIRA_PROJECT_PREFIX    Project prefix")
	fmt.Println("  ANTHROPIC_API_KEY      For Claude models")
	fmt.Println("  OPENAI_API_KEY         For GPT models")
	fmt.Println("  MCQ_PROMPTS_DIR        Custom templates directory")
	fmt.Println()
}

// listAllCommands displays all commands in a categorized list
func listAllCommands() {
	fmt.Println("📚 All MCQ Commands")
	fmt.Println("==================")
	fmt.Println()

	// Category definitions
	categories := map[string][]struct {
		name        string
		description string
		aliases     []string
	}{
		"🔧 Configuration": {
			{"config setup", "Interactive configuration setup", nil},
			{"config test", "Test current configuration", nil},
			{"config show", "Show current configuration", nil},
			{"context test", "Test repository context gathering", nil},
			{"config templates generate", "Generate example templates", nil},
			{"config templates validate", "Validate template syntax", nil},
			{"config templates list", "List available prompt types", nil},
		},
		"💻 Development Workflow": {
			{"build", "Build the binary", nil},
			{"build linux", "Cross-compile for Linux", nil},
			{"build windows", "Cross-compile for Windows", nil},
			{"test", "Run all tests", nil},
			{"test unit", "Run unit tests", nil},
			{"fmt", "Format code", nil},
			{"lint", "Run linters (golangci-lint, staticcheck)", []string{"-f", "--fix"}},
			{"staticcheck", "Run staticcheck", nil},
			{"reviewdog", "Run reviewdog", []string{"-p", "-s"}},
			{"vet", "Run go vet", nil},
			{"deps", "Manage dependencies (tidy, vendor)", nil},
			{"cover", "Generate test coverage report", nil},
			{"run", "Run main.go", nil},
			{"clean", "Prep for commit (fmt deps vet)", nil},
			{"setup", "Install dependencies", nil},
			{"install", "Install the binary", nil},
		},
		"📋 JIRA Integration": {
			{"jira show", "Display JIRA issue details", []string{"view", "display", "get"}},
			{"jira new", "Create JIRA issue from feature request", []string{"create", "add"}},
			{"jira new --dry-run", "Generate user story without creating ticket", nil},
		},
		"🐳 Docker": {
			{"docker build", "Build docker image", nil},
			{"docker run", "Run docker container", nil},
			{"docker push", "Push to registry", nil},
		},
		"📂 Git": {
			{"log", "Pretty git log", nil},
			{"gitclean", "Clean git working directory", nil},
		},
		"🚀 CI/CD": {
			{"ci", "Run CI checks", nil},
			{"all", "Run everything", nil},
		},
		"❓ Help & Info": {
			{"help", "Show help information", nil},
			{"examples", "Show usage examples", nil},
			{"workflow", "Show recommended workflows", nil},
			{"quick", "Show quick reference", nil},
			{"commands", "List all commands", nil},
			{"version", "Show version", []string{"v", "-v"}},
		},
		"📝 Templates": {
			{"templates generate", "Generate example template files", nil},
			{"templates validate", "Validate template syntax", nil},
			{"templates list", "List available prompt types", nil},
		},
	}

	// Print categories
	for catName, commands := range categories {
		fmt.Println(catName)
		fmt.Println(strings.Repeat("-", len(catName)-2))
		for _, cmd := range commands {
			aliasStr := ""
			if len(cmd.aliases) > 0 {
				aliases := make([]string, len(cmd.aliases))
				for i, alias := range cmd.aliases {
					if len(alias) == 1 {
						aliases[i] = "-" + alias
					} else {
						aliases[i] = "--" + alias
					}
				}
				aliasStr = " [" + strings.Join(aliases, ", ") + "]"
			}
			fmt.Printf("  %-30s %s%s\n", cmd.name, cmd.description, aliasStr)
		}
		fmt.Println()
	}

	fmt.Println("💡 Tip: Use 'mcq help <command>' for detailed information about any command")
	fmt.Println()
}
