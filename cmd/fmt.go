package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "-> go fmt",
	Long:  `This subcommand runs go fmt on all code`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Fmt()
	},
}

func init() {
	RootCmd.AddCommand(fmtCmd)
}
