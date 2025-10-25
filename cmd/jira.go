package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stevemcquaid/mcq/pkg/ai"
	"github.com/stevemcquaid/mcq/pkg/commands"
)

// jiraCmd represents the jira command
var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Jira integration commands",
	Long:  `Commands for interacting with Jira issues and projects.`,
}

// jiraShowCmd represents the jira show command
var jiraShowCmd = &cobra.Command{
	Use:     "show <issue-key>",
	Aliases: []string{"view", "display", "get"},
	Short:   "Display detailed information about a Jira issue",
	Long: `Show comprehensive details about a Jira issue including:
- Basic fields (title, summary, description, status, assignee, sprint, parent, related tickets)
- Comments (with author and timestamps) - you'll be prompted before displaying
- Work logs (coming soon)

If JIRA_PROJECT_PREFIX is set, you can use just the issue number:
  mcq jira show 123  (becomes PROJ-123 if prefix is "PROJ")

Examples:
  mcq jira show PROJ-123
  mcq jira show BUG-456
  mcq jira show 123  (requires JIRA_PROJECT_PREFIX=PROJ)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := args[0]
		commands.ShowJiraIssue(issueKey)
	},
}

// jiraNewCmd represents the jira new command
var jiraNewCmd = &cobra.Command{
	Use:     "new [flags] [--] <vague user story>",
	Aliases: []string{"create", "add"},
	Short:   "Create a new Jira issue from a vague user story using AI",
	Long: `Create a new Jira issue by converting a vague user story into a detailed user story using AI.

This command will:
1. Generate a detailed user story using AI (same as 'mcq ai jira')
2. Ask for confirmation before creating the Jira issue
3. Create a new Jira issue with the generated content as description
4. Copy the generated user story to clipboard

The issue will be created in the project specified by JIRA_PROJECT_PREFIX environment variable.

You can provide the user story in two ways:
1. Quoted: mcq jira new "Add dark mode to the application"
2. Unquoted with --: mcq jira new -- Add dark mode to the application

Examples:
  mcq jira new "Add dark mode to the application"
  mcq jira new -- Add dark mode to the application
  mcq jira new -v 8 "install/upgrade via homebrew"
  mcq jira new -v 8 -- install/upgrade via homebrew
  mcq jira new --model claude -- Improve user login process`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if -- is present in the arguments
		doubleDashIndex := -1
		for i, arg := range args {
			if arg == "--" {
				doubleDashIndex = i
				break
			}
		}

		var userStoryArgs []string
		if doubleDashIndex >= 0 {
			// -- syntax: everything after -- is the user story
			userStoryArgs = args[doubleDashIndex+1:]
			if len(userStoryArgs) == 0 {
				userStoryArgs = []string{""} // Empty user story
			}
		} else {
			// No -- found, treat all args as the user story (quoted syntax)
			userStoryArgs = args
		}

		// Get flags normally since we're not using DisableFlagParsing
		model, _ := cmd.Flags().GetString("model")
		verbosity, _ := cmd.Flags().GetInt("verbosity")
		autoDetect, _ := cmd.Flags().GetBool("auto-context")
		includeReadme, _ := cmd.Flags().GetBool("include-readme")
		includeGoMod, _ := cmd.Flags().GetBool("include-go-mod")
		includeCommits, _ := cmd.Flags().GetBool("include-commits")
		includeStructure, _ := cmd.Flags().GetBool("include-structure")
		includeConfigs, _ := cmd.Flags().GetBool("include-configs")
		maxCommits, _ := cmd.Flags().GetInt("max-commits")
		noContext, _ := cmd.Flags().GetBool("no-context")

		// Determine context configuration
		var contextConfig ai.ContextConfig
		if noContext {
			contextConfig = ai.ContextConfig{}
		} else if autoDetect || includeReadme || includeGoMod || includeCommits || includeStructure || includeConfigs {
			contextConfig = ai.ContextConfig{
				AutoDetect:       autoDetect,
				IncludeReadme:    includeReadme,
				IncludeGoMod:     includeGoMod,
				IncludeCommits:   includeCommits,
				IncludeStructure: includeStructure,
				IncludeConfigs:   includeConfigs,
				MaxCommits:       maxCommits,
				MaxFileSize:      50 * 1024,
			}
		} else {
			contextConfig = ai.PromptForContext()
		}

		if err := commands.JiraNew(userStoryArgs, model, verbosity, contextConfig); err != nil {
			// Error handling is done within JiraNew function
			// Exit with error code 1 to indicate failure
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(jiraCmd)
	jiraCmd.AddCommand(jiraShowCmd)
	jiraCmd.AddCommand(jiraNewCmd)

	// Jira configuration
	jiraCmd.PersistentFlags().String("url", "", "Jira instance URL (can also be set via JIRA_INSTANCE_URL env var)")
	jiraCmd.PersistentFlags().String("token", "", "Jira API token (can also be set via JIRA_API_TOKEN env var)")
	jiraCmd.PersistentFlags().String("username", "", "Jira username (for basic auth, can also be set via JIRA_USERNAME env var)")
	jiraCmd.PersistentFlags().String("password", "", "Jira password (for basic auth, can also be set via JIRA_PASSWORD env var)")
	jiraCmd.PersistentFlags().String("project-prefix", "", "Jira project prefix (can also be set via JIRA_PROJECT_PREFIX env var)")

	// AI flags for jira new command (for help display)
	jiraNewCmd.Flags().StringP("model", "m", "", "AI model to use: 'claude', 'gpt-4o', 'gpt-5', 'gpt-5-mini', or 'gpt-5-nano' (auto-detected if not specified)")
	jiraNewCmd.Flags().IntP("verbosity", "v", 0, "Set verbosity level: 0=off, 1=basic, 2=detailed, 3=verbose (use -v 8, not -v8)")
	jiraNewCmd.Flags().Bool("auto-context", false, "Automatically detect and include relevant repository context")
	jiraNewCmd.Flags().Bool("include-readme", false, "Include README content in context")
	jiraNewCmd.Flags().Bool("include-go-mod", false, "Include go.mod information in context")
	jiraNewCmd.Flags().Bool("include-commits", false, "Include recent commit messages in context")
	jiraNewCmd.Flags().Bool("include-structure", false, "Include directory structure in context")
	jiraNewCmd.Flags().Bool("include-configs", false, "Include configuration files in context")
	jiraNewCmd.Flags().Int("max-commits", 10, "Maximum number of recent commits to include")
	jiraNewCmd.Flags().Bool("no-context", false, "Skip context gathering entirely")

	// Bind flags to viper
	_ = viper.BindPFlag("jira.url", jiraCmd.PersistentFlags().Lookup("url"))
	_ = viper.BindPFlag("jira.token", jiraCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("jira.username", jiraCmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("jira.password", jiraCmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("jira.project_prefix", jiraCmd.PersistentFlags().Lookup("project-prefix"))

	// Set environment variable defaults
	_ = viper.BindEnv("jira.url", "JIRA_INSTANCE_URL")
	_ = viper.BindEnv("jira.token", "JIRA_API_TOKEN")
	_ = viper.BindEnv("jira.username", "JIRA_USERNAME")
	_ = viper.BindEnv("jira.password", "JIRA_PASSWORD")
	_ = viper.BindEnv("jira.project_prefix", "JIRA_PROJECT_PREFIX")
}
