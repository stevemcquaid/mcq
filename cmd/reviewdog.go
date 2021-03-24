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
		_ = commands.ReviewDog(PRNumFlag, SuggestionsFlag)
	},
}

var (
	PRNumFlag       int
	SuggestionsFlag bool
)

func init() {
	ReviewDogCmd.Flags().IntVarP(&PRNumFlag, "pull-request", "p", 0, "Comment lint comments in review")
	ReviewDogCmd.Flags().BoolVarP(&SuggestionsFlag, "suggest", "s", false, "Include suggested fixes in PR")
	RootCmd.AddCommand(ReviewDogCmd)
}
