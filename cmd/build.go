package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "-> go build",
	Long:  `This subcommand builds the binary`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.Build(viper.GetString("GIT_REPO"))
	},
}

var buildLinuxCmd = &cobra.Command{
	Use:   "linux",
	Short: "-> go build GOOS=linux",
	Long:  `This subcommand cross-compiles for linux`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.BuildLinux(viper.GetString("GIT_REPO") + ".linux.amd64.bin")
	},
}

var buildWindowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "-> go build GOOS=windows",
	Long:  `This subcommand cross-compiles for windows`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = commands.BuildWindows(viper.GetString("GIT_REPO") + ".windows.amd64.exe")
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(buildLinuxCmd)
	buildCmd.AddCommand(buildWindowsCmd)
}
