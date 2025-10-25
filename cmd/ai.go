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
	Long: `Convert a vague feature request into a detailed user story and copy it to clipboard.

The command will ask if you want to include repository context to improve the user story.
Use flags to skip the interactive prompts and specify context options directly.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		model, _ := cmd.Flags().GetString("model")
		verbosity, _ := cmd.Flags().GetInt("verbosity")

		// Get context flags
		autoDetect, _ := cmd.Flags().GetBool("auto-context")
		includeReadme, _ := cmd.Flags().GetBool("include-readme")
		includeGoMod, _ := cmd.Flags().GetBool("include-go-mod")
		includeCommits, _ := cmd.Flags().GetBool("include-commits")
		includeStructure, _ := cmd.Flags().GetBool("include-structure")
		includeConfigs, _ := cmd.Flags().GetBool("include-configs")
		maxCommits, _ := cmd.Flags().GetInt("max-commits")
		noContext, _ := cmd.Flags().GetBool("no-context")

		// Determine context configuration
		var contextConfig commands.ContextConfig
		if noContext {
			// Skip context gathering entirely
			contextConfig = commands.ContextConfig{}
		} else if autoDetect || includeReadme || includeGoMod || includeCommits || includeStructure || includeConfigs {
			// Use explicit flags
			contextConfig = commands.ContextConfig{
				AutoDetect:       autoDetect,
				IncludeReadme:    includeReadme,
				IncludeGoMod:     includeGoMod,
				IncludeCommits:   includeCommits,
				IncludeStructure: includeStructure,
				IncludeConfigs:   includeConfigs,
				MaxCommits:       maxCommits,
				MaxFileSize:      50 * 1024, // 50KB default
			}
		} else {
			// Ask user interactively
			contextConfig = commands.PromptForContext()
		}

		_ = commands.AIJira(args, model, verbosity, contextConfig)
	},
}

func init() {
	RootCmd.AddCommand(aiCmd)
	aiCmd.AddCommand(aiJiraCmd)

	// Add model flag
	aiJiraCmd.Flags().StringP("model", "m", "", "AI model to use: 'claude', 'gpt-4o', 'gpt-5', 'gpt-5-mini', or 'gpt-5-nano' (auto-detected if not specified)")

	// Add verbosity flag
	aiJiraCmd.Flags().IntP("verbosity", "v", 0, "Set verbosity level: 0=off, 1=basic, 2=detailed, 3=verbose")

	// Add context flags
	aiJiraCmd.Flags().Bool("auto-context", false, "Automatically detect and include relevant repository context")
	aiJiraCmd.Flags().Bool("include-readme", false, "Include README content in context")
	aiJiraCmd.Flags().Bool("include-go-mod", false, "Include go.mod information in context")
	aiJiraCmd.Flags().Bool("include-commits", false, "Include recent commit messages in context")
	aiJiraCmd.Flags().Bool("include-structure", false, "Include directory structure in context")
	aiJiraCmd.Flags().Bool("include-configs", false, "Include configuration files in context")
	aiJiraCmd.Flags().Int("max-commits", 10, "Maximum number of recent commits to include")
	aiJiraCmd.Flags().Bool("no-context", false, "Skip context gathering entirely")
}
