package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stevemcquaid/mcq/pkg/commands"
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI-powered commands",
	Long:  `AI-powered commands for development tasks`,
}

var aiJiraCmd = &cobra.Command{
	Use:   "jira [feature request]",
	Short: "Convert vague feature request to user story",
	Long:  `Convert a vague feature request into a detailed user story and copy it to clipboard`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		model, _ := cmd.Flags().GetString("model")
		verbosity, _ := cmd.Flags().GetInt("verbosity")
		_ = commands.AIJira(args, model, verbosity)
	},
}

func init() {
	RootCmd.AddCommand(aiCmd)
	aiCmd.AddCommand(aiJiraCmd)

	// Add model flag
	aiJiraCmd.Flags().StringP("model", "m", "", "AI model to use: 'claude', 'gpt-4o', 'gpt-5', 'gpt-5-mini', or 'gpt-5-nano' (auto-detected if not specified)")

	// Add verbosity flag
	aiJiraCmd.Flags().IntP("verbosity", "v", 0, "Set verbosity level: 0=off, 1=basic, 2=detailed, 3=verbose")
}
