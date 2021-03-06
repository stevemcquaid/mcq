package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

// ## Prep for commit - run make fmt, vendor, tidy
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "-> fmt deps vet",
	Long:  `This subcommand preps for commit: runs fmt, fumpt, deps & vet`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Clean()
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
