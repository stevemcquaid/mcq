package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Install() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go install",
				Function: shell.PrettyRun,
			},
		},
	)
}
