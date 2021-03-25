package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.2"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Version",
	Aliases: []string{"v", "-v"},
	Long:    `This subcommand returns the version of the CLI utility`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
