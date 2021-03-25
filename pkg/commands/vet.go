package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Vet() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go vet $(go list ./... | grep -v vendor) | grep -v '.pb.go:' | tee /dev/stderr",
				Function: shell.PrettyRun,
			},
		},
	)
}
