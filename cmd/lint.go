package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var LintCmd = &cobra.Command{
	Use:   "lint",
	Short: "-> golangci-lint, staticcheck",
	Long:  `This subcommand runs static analysis tools`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Lint()
	},
}

func init() {
	RootCmd.AddCommand(LintCmd)
}
