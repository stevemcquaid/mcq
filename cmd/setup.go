package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "install dependencies",
	Long:  `This subcommand installs build and lint dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Setup()
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
