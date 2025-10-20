package commands

import (
	"github.com/stevemcquaid/mcq/pkg/shell"
)

// Install all the build and lint dependencies
func Setup() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			//&shell.StringFunction{
			//	Arg:      "go install go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0",
			//	Function: shell.PrettyRun,
			//},
			&shell.StringFunction{
				Arg:      "go install honnef.co/go/tools/cmd/staticcheck@latest",
				Function: shell.PrettyRun,
			},

			// &shell.StringFunction{
			// 	Arg:      "go get -u github.com/pierrre/gotestcover",
			// 	Function: shell.PrettyRun,
			// },
			//
			// &shell.StringFunction{
			// 	Arg:      "go get -u golang.org/x/tools/cmd/cover",
			// 	Function: shell.PrettyRun,
			// },

			&shell.StringFunction{
				Arg:      "go install golang.org/x/tools/cmd/goimports@latest",
				Function: shell.PrettyRun,
			},

			// &shell.StringFunction{
			// 	Arg:      "go get -u github.com/reviewdog/reviewdog/cmd/reviewdog",
			// 	Function: shell.PrettyRun,
			// },

			&shell.StringFunction{
				Arg:      "go install mvdan.cc/gofumpt@latest",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go install github.com/wadey/gocovmerge@latest",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go install github.com/uudashr/gocognit/cmd/gocognit@latest",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "go install github.com/fzipp/gocyclo/cmd/gocyclo@latest",
				Function: shell.PrettyRun,
			},

			&shell.StringFunction{
				Arg:      "mkdir -p build",
				Function: shell.PrettyRun,
			},
		},
	)
}
