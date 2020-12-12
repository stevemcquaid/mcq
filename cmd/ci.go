package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var CICmd = &cobra.Command{
	Use:   "ci",
	Short: "Run almost everything",
	Long:  `This subcommand runs all the tests and code checks`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.CI()
	},
}

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "Run everything",
	Long:  `This subcommand runs everything`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.All()
	},
}

func init() {
	RootCmd.AddCommand(CICmd)
	RootCmd.AddCommand(AllCmd)
}
