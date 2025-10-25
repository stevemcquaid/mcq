package commands

import (
	"fmt"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

// Run all the tests and code checks
func All() error {
	dockerImage, err := GetDockerImage()
	if err != nil {
		fmt.Printf("\n \n \n GOT ERROR: %s", err)
		return err
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Setup,
			},
			&shell.VoidFunction{
				Function: CI,
			},
			&shell.StringFunction{
				Arg:      dockerImage,
				Function: DockerBuild,
			},
		},
	)
}

// Run all the tests and code checks
func CI() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Deps,
			},
			&shell.VoidFunction{
				Function: Fmt,
			},
			&shell.VoidFunction{
				Function: Vet,
			},
			&shell.BoolFunction{
				Function: Lint,
				Arg:      false,
			},
			&shell.VoidFunction{
				Function: Test,
			},
			//&shell.StringSliceFunction{
			//	Function: Install,
			//	Arg:      []string{},
			//},
		},
	)
}
