package commands

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

// Run all linters
func Lint(fixFlag bool) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.BoolFunction{
				Arg:      fixFlag,
				Function: GolangCI,
			},
			&shell.VoidFunction{
				Function: StaticCheck,
			},
		},
	)
}

var StaticCheckCommand = []string{
	"staticcheck -fail -tests -checks=\"all,-ST1000,-ST1021,-ST1020\" ./...",
}

func StaticCheck() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      strings.Join(StaticCheckCommand, " "),
				Function: shell.PrettyRun,
			},
		},
	)
}

var GolangciLintCommand = []string{
	"golangci-lint run",
	"--deadline=30m",
	"--disable-all",
	"--no-config",
	"--issues-exit-code=1",
	"--enable=bodyclose",
	"--enable=deadcode",
	"--enable=dupl",
	"--enable=errcheck",
	"--enable=gocognit",
	"--enable=goconst ",
	"--enable=gocyclo",
	"--enable=gofmt",
	"--enable=gofumpt",
	"--enable=goimports",
	"--enable=golint",
	"--enable=gomodguard",
	"--enable=gosec ",
	"--enable=govet",
	"--enable=ineffassign",
	"--enable=interfacer ",
	"--enable=megacheck",
	"--enable=misspell",
	"--enable=nakedret",
	"--enable=prealloc",
	"--enable=rowserrcheck",
	"--enable=staticcheck",
	"--enable=structcheck ",
	"--enable=stylecheck",
	"--enable=typecheck",
	"--enable=unconvert ",
	"--enable=unparam",
	"--enable=varcheck",
	"--enable=whitespace",
}

func getGolangCICommandWithFix(fix bool) string {
	var command []string
	if fix {
		command = append(GolangciLintCommand, "--fix")
	} else {
		command = GolangciLintCommand
	}

	return strings.Join(command, " ")
}

func GolangCI(fix bool) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      getGolangCICommandWithFix(fix),
				Function: shell.PrettyRun,
			},
		},
	)
}

func ReviewDog(pr int, suggest bool) error {
	gitOrg, gitRepo, err := GetModules()
	if err != nil {
		return err

	}

	// dont include suggestions
	lintCmd:= getGolangCICommandWithFix(false)
	command:= []string{
		fmt.Sprintf("export CI_PULL_REQUEST=%d;", pr),
		fmt.Sprintf("export CI_REPO_OWNER=%s;", gitOrg),
		fmt.Sprintf("export CI_REPO_NAME=%s;", gitRepo),
		"export CI_COMMIT=$(git rev-parse HEAD);",
		lintCmd + " | reviewdog -f=golangci-lint -diff=\"git diff FETCH_HEAD\" -reporter=github-pr-review",
	}

	if suggest {
		// include suggestions
		command = []string{
			fmt.Sprintf("export CI_PULL_REQUEST=%d;", pr),
			fmt.Sprintf("export CI_REPO_OWNER=%s;", gitOrg),
			fmt.Sprintf("export CI_REPO_NAME=%s;", gitRepo),
			"export CI_COMMIT=$(git rev-parse HEAD);",
			"export TMPFILELINT=$(mktemp);",
			getGolangCICommandWithFix(false) + " --out-format=json > $TMPFILELINT;",
			"export TMPFILEDIFF=$(mktemp);",
			getGolangCICommandWithFix(true) + ";",
			"git diff > $TMPFILEDIFF;",
			"git stash -u && git stash drop;",
			"reviewdog -f=diff -f.diff.strip=1 -reporter=github-pr-review < \"${TMPFILEDIFF}\"",
			// "cat $TMPFILELINT | reviewdog -f=golangci-lint -f.diff.strip=1 -diff=\"cat $TMPFILEDIFF\" -reporter=github-pr-review;",
			// "cat $TMPFILELINT | reviewdog -name=\"gofmt\" -f=golangci-lint -diff=\"cat $TMPFILEDIFF\" -reporter=github-pr-review;",
		}


		
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      strings.Join(command, ""),
				Function: shell.PrettyRun,
			},
		},
	)
}
