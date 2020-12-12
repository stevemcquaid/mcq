package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stevemcquaid/mcq/pkg/commands"
)

// dockerRunCmd represents a docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "docker build, run, push",
	Long:  `Various docker tasks`,
}

// dockerRunCmd represents a docker command
var dockerBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "docker build",
	Long:  `This subcommand builds the dockerfile`,
	Run: func(cmd *cobra.Command, args []string) {
		gitOrg := viper.GetString("GIT_ORG")
		gitRepo := viper.GetString("GIT_REPO")
		dockerBase := path.Join(gitOrg, gitRepo)
		dockerImage := fmt.Sprintf("%s:%s", dockerBase, "latest")
		commands.DockerBuild(dockerImage)
	},
}

// dockerRunCmd represents a docker command
var dockerRunCmd = &cobra.Command{
	Use:   "run",
	Short: "docker run",
	Long:  `This subcommand runs docker`,
	Run: func(cmd *cobra.Command, args []string) {
		gitOrg := viper.GetString("GIT_ORG")
		gitRepo := viper.GetString("GIT_REPO")
		dockerBase := path.Join(gitOrg, gitRepo)
		dockerImage := fmt.Sprintf("%s:%s", dockerBase, "latest")
		commands.DockerRun(dockerImage)
	},
}

// dockerPushCmd represents a docker command
var dockerPushCmd = &cobra.Command{
	Use:   "push",
	Short: "docker push",
	Long:  `This subcommand runs docker push`,
	Run: func(cmd *cobra.Command, args []string) {
		gitOrg := viper.GetString("GIT_ORG")
		gitRepo := viper.GetString("GIT_REPO")
		dockerBase := path.Join(gitOrg, gitRepo)
		dockerImage := fmt.Sprintf("%s:%s", dockerBase, "latest")
		commands.DockerPush(dockerImage)
	},
}

func init() {
	RootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(dockerRunCmd)
	dockerCmd.AddCommand(dockerBuildCmd)
	dockerCmd.AddCommand(dockerPushCmd)
}
