package commands

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

// builds the binary
func Build(binaryName string, args []string) error {
	filePath := strings.Join(args, " ")

	if len(filePath) == 0 {
		filePath = "./"
	}
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("go build -o bin/%s %s", binaryName, filePath),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func BuildLinux(binaryName string, args []string) error {
	filePath := strings.Join(args, " ")

	if len(filePath) == 0 {
		filePath = "./"
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("GOOS=linux GOARCH=amd64 go build -o bin/%s %s", binaryName, filePath),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func BuildWindows(binaryName string, args []string) error {
	filePath := strings.Join(args, " ")

	if len(filePath) == 0 {
		filePath = "./"
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("GOOS=windows GOARCH=amd64 go build -o bin/%s %s", binaryName, filePath),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}
