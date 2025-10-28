package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
1. Generate a detailed user story using AI
2. Display the generated user story for review
3. Ask for confirmation before creating the Jira issue (unless --dry-run is used)
4. Create a new Jira issue with the generated content as description
5. Copy the generated user story to clipboard

The issue will be created in the project specified by JIRA_PROJECT_PREFIX environment variable.

You can provide the user story in two ways:
1. Quoted (recommended for simple text):
   mcq jira new "Add dark mode to the application"
   
2. Unquoted with -- flag for unformatted input (useful for complex text or when you don't want to quote):
   mcq jira new -- Add dark mode to the application
   
The -- flag tells the command to treat everything after it as the user story text, useful when your text contains quotes, special characters, or you prefer not to quote it.

You can use --dry-run to generate the user story without creating the JIRA issue:
  mcq jira new --dry-run "Add dark mode"
  
This is useful when you want to see what would be generated without actually creating a ticket.

Examples:
  mcq jira new "Add dark mode to the application"
  mcq jira new -- Add dark mode to the application
  mcq jira new --dry-run "Add dark mode"  # Generate without creating ticket
  mcq jira new -v 8 -- install/upgrade via homebrew
  mcq jira new --model claude -- Improve user login process`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Cobra automatically handles -- as a flag terminator
		// Everything after -- is passed as positional args
		// So args here already contains only the user story text

		// Join all args into a single user story string
		userStoryArgs := args
		if len(args) == 0 {
			userStoryArgs = []string{""}
		}

		// Get flags
		model, _ := cmd.Flags().GetString("model")
		verbosity, _ := cmd.Flags().GetInt("verbosity")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Extract context configuration
		contextConfig := extractContextConfig(cmd)

		if err := commands.JiraNew(userStoryArgs, model, verbosity, contextConfig, dryRun); err != nil {
			// Error handling is done within JiraNew function
			// Exit with error code 1 to indicate failure
			os.Exit(1)
		}
	},
}

// jiraUpdateCmd represents the jira update command
var jiraUpdateCmd = &cobra.Command{
	Use:     "update <issue-key>",
	Aliases: []string{"improve", "upgrade"},
	Short:   "Improve and update a Jira issue description using AI",
	Long: `Improve an existing Jira issue description using AI and optionally update it.

This command will:
1. Fetch the current issue description
2. Use AI to improve and enhance it
3. Show a side-by-side comparison of BEFORE/AFTER
4. Copy the improved description to clipboard
5. Ask for confirmation before updating (unless --dry-run is used)
6. Update the Jira issue with the improved description

You can use --dry-run to generate the improved description without updating the issue:
  mcq jira update --dry-run RS-1812
  
This is useful when you want to see what would be generated without actually updating.

Examples:
  mcq jira update RS-1812
  mcq jira update --dry-run RS-1812
  mcq jira update --model claude RS-1812
  mcq jira update --no-context RS-1812`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := args[0]

		// Get flags
		model, _ := cmd.Flags().GetString("model")
		verbosity, _ := cmd.Flags().GetInt("verbosity")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Extract context configuration
		contextConfig := extractContextConfig(cmd)

		if err := commands.JiraUpdate(issueKey, model, verbosity, contextConfig, dryRun); err != nil {
			// Error handling is done within JiraUpdate function
			// Exit with error code 1 to indicate failure
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(jiraCmd)
	jiraCmd.AddCommand(jiraShowCmd)
	jiraCmd.AddCommand(jiraNewCmd)
	jiraCmd.AddCommand(jiraUpdateCmd)

	// Jira configuration
	jiraCmd.PersistentFlags().String("url", "", "Jira instance URL (can also be set via JIRA_INSTANCE_URL env var)")
	jiraCmd.PersistentFlags().String("token", "", "Jira API token (can also be set via JIRA_API_TOKEN env var)")
	jiraCmd.PersistentFlags().String("username", "", "Jira username (for basic auth, can also be set via JIRA_USERNAME env var)")
	jiraCmd.PersistentFlags().String("password", "", "Jira password (for basic auth, can also be set via JIRA_PASSWORD env var)")
	jiraCmd.PersistentFlags().String("project-prefix", "", "Jira project prefix (can also be set via JIRA_PROJECT_PREFIX env var)")

	// AI flags for jira new and update commands
	addAIFlags(jiraNewCmd)
	addAIFlags(jiraUpdateCmd)

	// Add dry-run flags
	jiraNewCmd.Flags().Bool("dry-run", false, "Generate user story without creating JIRA issue (alias for 'mcq ai jira')")
	jiraUpdateCmd.Flags().Bool("dry-run", false, "Generate improved description without updating JIRA issue")

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
