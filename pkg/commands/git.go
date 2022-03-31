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

func GitClean() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "git reset --hard HEAD",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "git clean -fd",
				Function: shell.PrettyRun,
			},
		})
}
