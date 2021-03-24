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
	"--sort-results",
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

func GolangCI(fix bool) error {
	var command []string
	if fix {
		command = append(GolangciLintCommand, "--fix")
	} else {
		command = GolangciLintCommand
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      strings.Join(command, " "),
				Function: shell.PrettyRun,
			},
		},
	)
}

func ReviewDog(PR int) error {
	gitOrg, gitRepo, err := GetModules()
	if err != nil {
		return err
	}

	// lintCmd := "golint ./..."
	lintCmd := "staticcheck -fail -tests -checks=\"all,-ST1000,-ST1021,-ST1020\" ./..."

	command := []string{
		fmt.Sprintf("export CI_PULL_REQUEST=%d;", PR),
		fmt.Sprintf("export CI_REPO_OWNER=%s;", gitOrg),
		fmt.Sprintf("export CI_REPO_NAME=%s;", gitRepo),
		"export CI_COMMIT=$(git rev-parse HEAD);",
		lintCmd + " | reviewdog -f=golint -diff=\"git diff FETCH_HEAD\" -reporter=github-pr-review",
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
