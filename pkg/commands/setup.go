package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

// Install all the build and lint dependencies
func Setup() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "brew install golangci-lint",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go get honnef.co/go/tools/cmd/staticcheck@2020.1.6",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go install honnef.co/go/tools/cmd/staticcheck@2020.1.6",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go get -u github.com/pierrre/gotestcover",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go get -u golang.org/x/tools/cmd/cover",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go get -u golang.org/x/tools/cmd/goimports",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go install golang.org/x/tools/cmd/goimports",
				Function: shell.PrettyRun,
			},

			// &shell.StringFunction{
			// 	Arg:      "go get -u github.com/reviewdog/reviewdog/cmd/reviewdog",
			// 	Function: shell.PrettyRun,
			// },

			&shell.StringFunction{
				Arg:      "go get -u mvdan.cc/gofumpt",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go install mvdan.cc/gofumpt",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "mkdir -p build",
				Function: shell.PrettyRun,
			},
		},
	)
}
