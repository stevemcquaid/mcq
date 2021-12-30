package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Test() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Deps,
			},
			&shell.VoidFunction{
				Function: TestUnit,
			},
		},
	)
}

func TestUnit() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go test -cover -covermode=atomic -coverprofile=build/unit.out $(go list ./... | grep -v /vendor/)  -run .",
				Function: shell.PrettyRun,
			},
		},
	)
}

func TestIntegration() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go test -cover -covermode=atomic -coverprofile=build/integration.out -tags=integration $(go list ./... | grep -v /vendor/) -run $(TEST_PATTERN)",
				Function: shell.PrettyRun,
			},
		},
	)
}
