package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Fmt() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s -l \"$file\"; done",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w -l \"$file\"; done",
				Function: shell.PrettyRun,
			},
		},
	)
}
