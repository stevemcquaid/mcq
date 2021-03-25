package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Log() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "git log --all --decorate --oneline --graph",
				Function: shell.PrettyRun,
			},
		})
}
