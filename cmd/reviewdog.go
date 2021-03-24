package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var ReviewDogCmd = &cobra.Command{
	Use:   "reviewdog",
	Short: "-> reviewdog",
	Long:  `Runs reviewdog`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.ReviewDog(PRNumFlag)
	},
}

var PRNumFlag int

func init() {
	ReviewDogCmd.Flags().IntVarP(&PRNumFlag, "PullRequest", "p", 0, "Comment lint comments in review")
	RootCmd.AddCommand(ReviewDogCmd)
}
