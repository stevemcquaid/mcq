package commands

import (
	"fmt"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

func DockerBuild(dockerImage string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker build --target final -t %s .", dockerImage),
				Function: shell.PrettyRun,
			},
		},
	)
}

// @TODO - figure out port requirements
func DockerRun(dockerImage string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker build --target final -t %s .", dockerImage),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      dockerImage,
				Function: DockerBuild,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker run -it -P %s .", dockerImage),
				Function: shell.PrettyRun,
			},
		},
	)
}

// @TODO - figure out port requirements
func DockerPush(dockerImage string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      dockerImage,
				Function: DockerBuild,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker push %s", dockerImage),
				Function: shell.PrettyRun,
			},
		},
	)
}
