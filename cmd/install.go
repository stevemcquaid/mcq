package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "-> go install",
	Long:  `This subcommand installs the binary into gopath`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Install()
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
