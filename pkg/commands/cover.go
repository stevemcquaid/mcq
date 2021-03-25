package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

// Run all the tests and opens the coverage report
func Cover() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Test,
			},
			&shell.StringFunction{
				Arg:      "gocovmerge build/unit.out > build/all.out",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go tool cover -html=build/all.out",
				Function: shell.PrettyRun,
			},
		},
	)
}
