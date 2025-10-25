package cmd

import (
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
	Use:   "show <issue-key>",
	Short: "Display detailed information about a Jira issue",
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

func init() {
	RootCmd.AddCommand(jiraCmd)
	jiraCmd.AddCommand(jiraShowCmd)

	// Jira configuration
	jiraCmd.PersistentFlags().String("url", "", "Jira instance URL (can also be set via JIRA_INSTANCE_URL env var)")
	jiraCmd.PersistentFlags().String("token", "", "Jira API token (can also be set via JIRA_API_TOKEN env var)")
	jiraCmd.PersistentFlags().String("username", "", "Jira username (for basic auth, can also be set via JIRA_USERNAME env var)")
	jiraCmd.PersistentFlags().String("password", "", "Jira password (for basic auth, can also be set via JIRA_PASSWORD env var)")
	jiraCmd.PersistentFlags().String("project-prefix", "", "Jira project prefix (can also be set via JIRA_PROJECT_PREFIX env var)")

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
