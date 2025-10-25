package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mcq",
	Short: "A Makefile helper with JIRA and AI integration",
	Long: `MCQ is a powerful CLI tool that helps streamline your development workflow
with JIRA integration and AI-powered user story generation.

Key Features:
• JIRA issue management and creation
• AI-powered user story generation
• Interactive configuration setup
• Context-aware content generation

Quick Start:
1. Run 'mcq config setup' to configure JIRA and AI settings
2. Run 'mcq config test' to verify your configuration
3. Try 'mcq ai jira "Add dark mode"' to generate a user story
4. Use 'mcq jira new "Add dark mode"' to create a JIRA issue

For more information, run 'mcq help' or 'mcq examples'.`,
	Run: func(cmd *cobra.Command, args []string) {
		showMainHelp()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Load the PWD golang module name
	gitOrg, gitRepo, err := commands.GetModules()
	if err != nil {
		fmt.Println("unable to set GIT_ORG + GIT_REPO")
	}

	viper.Set("GIT_ORG", gitOrg)
	viper.Set("GIT_REPO", gitRepo)
}

// showMainHelp displays the main help when no arguments are provided
func showMainHelp() {
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
	fmt.Println("  jira create      Alias for 'jira new'")
	fmt.Println("  jira add         Alias for 'jira new'")
	fmt.Println()
	fmt.Println("🤖 AI Commands:")
	fmt.Println("  ai jira <story>  Convert feature request to user story")
	fmt.Println("  ai story         Alias for 'ai jira'")
	fmt.Println("  ai generate      Alias for 'ai jira'")
	fmt.Println()
	fmt.Println("❓ Help Commands:")
	fmt.Println("  help             Show comprehensive help information")
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

	// More Information
	fmt.Println("📚 More Information:")
	fmt.Println("-------------------")
	fmt.Println("• Run 'mcq help' for comprehensive help")
	fmt.Println("• Run 'mcq examples' for detailed usage examples")
	fmt.Println("• Run 'mcq workflow' for recommended workflows")
	fmt.Println("• Run 'mcq <command> --help' for command-specific help")
	fmt.Println()
}
