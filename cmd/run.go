package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "-> go run main.go",
	Long:  `This subcommand runs the code`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Run()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
