package commands

import (
	"fmt"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

// builds the binary
func Build(binaryName string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("go build -o bin/%s ./", binaryName),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func BuildLinux(binaryName string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("GOOS=linux GOARCH=amd64 go build -o bin/%s ./", binaryName),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func BuildWindows(binaryName string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("GOOS=windows GOARCH=amd64 go build -o bin/%s ./", binaryName),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}
