package commands

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

// Run all linters
func Lint(fixFlag bool) error {
	return shell.OrderedRunnerIgnoreErrors(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: StaticCheck,
			},
			&shell.BoolFunction{
				Arg:      fixFlag,
				Function: GolangCI,
			},
			&shell.BoolFunction{
				Function: GoCyclo,
			},
			&shell.BoolFunction{
				Function: GoCogintCommand,
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
	"--timeout=1m",
	//"--disable-all",
	//"--no-config",
	//"--issues-exit-code=1",
	//"--enable=bodyclose",
	//"--enable=dupl",
	//"--enable=errcheck",
	//"--enable=gocognit",
	//"--enable=goconst ",
	//"--enable=gocyclo",
	//"--enable=gofmt",
	//"--enable=gofumpt",
	//"--enable=goimports",
	//"--enable=gomodguard",
	//"--enable=gosec ",
	//"--enable=govet",
	//"--enable=ineffassign",
	//"--enable=megacheck",
	//"--enable=misspell",
	//"--enable=nakedret",
	//"--enable=prealloc",
	//"--enable=revive",
	//"--enable=rowserrcheck",
	//"--enable=staticcheck",
	//"--enable=stylecheck",
	//"--enable=typecheck",
	//"--enable=unconvert ",
	//"--enable=unparam",
	//"--enable=unused",
	//"--enable=whitespace",
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

var GocogintCommand = []string{
	"gocognit -over 10  -ignore \"_test|testdata|vendor/*\" .",
}

func GoCogintCommand(fix bool) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      strings.Join(GocogintCommand, " "),
				Function: shell.PrettyRun,
			},
		},
	)
}

var GocycloCommand = []string{
	"gocyclo -over 25 .",
}

func GoCyclo(fix bool) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      strings.Join(GocycloCommand, " "),
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
	lintCmd := getGolangCICommandWithFix(false)
	command := []string{
		fmt.Sprintf("export CI_PULL_REQUEST=%d;", pr),
		fmt.Sprintf("export CI_REPO_OWNER=%s;", gitOrg),
		fmt.Sprintf("export CI_REPO_NAME=%s;", gitRepo),
		"export CI_COMMIT=$(git rev-parse HEAD);",
		lintCmd + " --out-format=line-number | reviewdog -name=\"golangci-lint\" -f=golangci-lint -diff=\"git diff FETCH_HEAD\" -reporter=github-pr-review",
	}

	if suggest {
		// include suggestions
		command = []string{
			fmt.Sprintf("export CI_PULL_REQUEST=%d;", pr),
			fmt.Sprintf("export CI_REPO_OWNER=%s;", gitOrg),
			fmt.Sprintf("export CI_REPO_NAME=%s;", gitRepo),
			"export CI_COMMIT=$(git rev-parse HEAD);",
			"export TMPFILEDIFF=$(mktemp);",
			getGolangCICommandWithFix(true) + " --out-format=line-number; ",
			"git diff > $TMPFILEDIFF;",
			"git stash -u && git stash drop;",
			"reviewdog -name=\"golangci-lint\" -f=diff -f.diff.strip=1 -reporter=github-pr-review < \"${TMPFILEDIFF}\"",
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
