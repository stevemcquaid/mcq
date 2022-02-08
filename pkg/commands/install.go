package commands

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

func Install(args []string) error {
	filePath := strings.Join(args, " ")

	if len(filePath) == 0 {
		filePath = "./"
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("go install %s", filePath),
				Function: shell.PrettyRun,
			},
		},
	)
}
