package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "-> go mod tidy, download, vendor",
	Long:  `This subcommand runs go mod tidy, download & vendor `,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Deps()
	},
}

func init() {
	RootCmd.AddCommand(depsCmd)
}
