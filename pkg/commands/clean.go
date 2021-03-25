package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Clean() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Fmt,
			},
			&shell.VoidFunction{
				Function: Fumpt,
			},
			&shell.StringFunction{
				Arg:      "go mod tidy",
				Function: shell.PrettyRun,
			},
			&shell.VoidFunction{
				Function: Deps,
			},
			&shell.VoidFunction{
				Function: Vet,
			},
		},
	)
}
