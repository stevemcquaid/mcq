package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Deps() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go mod tidy",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go mod download",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go mod vendor",
				Function: shell.PrettyRun,
			},
		},
	)
}
