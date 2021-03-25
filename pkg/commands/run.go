package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Run() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go run main.go",
				Function: shell.PrettyRun,
			},
		},
	)
}
