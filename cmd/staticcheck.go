package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var StaticCheckCmd = &cobra.Command{
	Use:   "staticcheck",
	Short: "-> staticcheck",
	Long:  `This subcommand runs staticcheck static-analysis tools`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.StaticCheck()
	},
}

func init() {
	RootCmd.AddCommand(StaticCheckCmd)
}
