package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "-> ~git log --graph --oneline --decorate --all",
	Long:  `This subcommand prettyPrints the git log`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Log()
	},
}

func init() {
	RootCmd.AddCommand(LogCmd)
}
