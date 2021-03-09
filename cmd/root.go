package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mcq",
	Short: "A Makefile helper",
	Long:  `This application provides shortcuts to common development tasks`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("make all")
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Load the PWD golang module name
	gitOrg, gitRepo, err := commands.GetModules()
	if err != nil {
		fmt.Println("unable to set GIT_ORG + GIT_REPO")
	}

	viper.Set("GIT_ORG", gitOrg)
	viper.Set("GIT_REPO", gitRepo)
}
