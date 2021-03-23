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
		_ = commands.Lint(FixFlag)
	},
}

var FixFlag bool

func init() {
	LintCmd.Flags().BoolVarP(&FixFlag, "fix", "f", false, "Fix found issues (if it's supported by the linter)")
	RootCmd.AddCommand(LintCmd)
}
