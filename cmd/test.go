package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var testUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "-> go test -tags=unit",
	Long:  `This subcommand runs unit tests`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.TestUnit()
	},
}

var testIntegratinoCmd = &cobra.Command{
	Use:   "unit",
	Short: "-> go test -tags=integration",
	Long:  `This subcommand runs intregation tests`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.TestIntegration()
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "-> go test",
	Long:  `This subcommand runs all tests`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Test()
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testUnitCmd)
	testCmd.AddCommand(testIntegratinoCmd)
}
