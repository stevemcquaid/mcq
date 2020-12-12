package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var coverCmd = &cobra.Command{
	Use:   "cover",
	Short: "-> go tool cover",
	Long:  `This subcommand runs all the tests and opens the coverage report`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Cover()
	},
}

func init() {
	RootCmd.AddCommand(coverCmd)
}
