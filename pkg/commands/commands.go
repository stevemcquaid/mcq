package commands

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

func DockerBuild(dockerImage string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker build --target final -t %s .", dockerImage),
				Function: shell.PrettyRun,
			},
		},
	)
}

// @TODO - figure out port requirements
func DockerRun(dockerImage string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker build --target final -t %s .", dockerImage),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      dockerImage,
				Function: DockerBuild,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker run -it -P %s .", dockerImage),
				Function: shell.PrettyRun,
			},
		},
	)
}

// @TODO - figure out port requirements
func DockerPush(dockerImage string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      dockerImage,
				Function: DockerBuild,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("docker push %s", dockerImage),
				Function: shell.PrettyRun,
			},
		},
	)
}

// builds the binary
func Build(binaryName string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("go build -o bin/%s ./", binaryName),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func BuildLinux(binaryName string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("GOOS=linux GOARCH=amd64 go build -o bin/%s ./", binaryName),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func BuildWindows(binaryName string) error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      fmt.Sprintf("GOOS=windows GOARCH=386 go build -o bin/%s ./", binaryName),
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      fmt.Sprintf("chmod u+x bin/%s", binaryName),
				Function: shell.PrettyRun,
			},
		},
	)
}

func Run() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go run main.go",
				Function: shell.PrettyRun,
			},
		},
	)
}

func Install() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go install",
				Function: shell.PrettyRun,
			},
		},
	)
}

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

func Clean() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Fmt,
			},
			&shell.StringFunction{
				Arg:      "go mod tidy",
				Function: shell.PrettyRun,
			},
			&shell.VoidFunction{
				Function: Deps,
			},
			&shell.VoidFunction{
				Function: Vet,
			},
		},
	)
}

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

func Deps() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go mod tidy",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go mod download",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "go mod vendor",
				Function: shell.PrettyRun,
			},
		},
	)
}

// Install all the build and lint dependencies
func Setup() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.32.2",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck@2020.1.6",
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
				Arg:      "go get -u github.com/reviewdog/reviewdog/cmd/reviewdog",
				Function: shell.PrettyRun,
			},
			&shell.StringFunction{
				Arg:      "mkdir -p build",
				Function: shell.PrettyRun,
			},
		},
	)
}

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
				Arg:      "go test -cover -covermode=atomic -coverprofile=build/unit.out $(go list ./... | grep -v /vendor/)  -run . -timeout=2m",
				Function: shell.PrettyRun,
			},
		},
	)
}

func TestIntegration() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "go test -cover -covermode=atomic -coverprofile=build/integration.out -tags=integration $(go list ./... | grep -v /vendor/) -run $(TEST_PATTERN) -timeout=2m",
				Function: shell.PrettyRun,
			},
		},
	)
}

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

func StaticCheck() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "staticcheck -fail -tests -checks=\"all,-ST1000,-ST1021,-ST1020\" ./...",
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
	// userName, err := GetUserName()
	// if err != nil {
	// 	return err
	// }

	gitOrg, gitRepo, err := GetModules()
	if err != nil {
		return err
	}

	command := []string{
		fmt.Sprintf("export CI_PULL_REQUEST=%d;", PR),
		fmt.Sprintf("export CI_REPO_OWNER=%s;", gitOrg),
		fmt.Sprintf("export CI_REPO_NAME=%s;", gitRepo),
		fmt.Sprintf("export CI_COMMIT=$(git rev-parse HEAD);"),
		"golint ./... | reviewdog -f=golint -diff=\"git diff FETCH_HEAD\" -reporter=github-pr-review",
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

// Run all the tests and code checks
func All() error {
	dockerImage, err := GetDockerImage()
	if err != nil {
		fmt.Printf("\n \n \n GOT ERROR: %s", err)
		return err
	}

	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Setup,
			},
			&shell.VoidFunction{
				Function: CI,
			},
			&shell.StringFunction{
				Arg:      dockerImage,
				Function: DockerBuild,
			},
		},
	)
}

// Run all the tests and code checks
func CI() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.VoidFunction{
				Function: Deps,
			},
			&shell.VoidFunction{
				Function: Fmt,
			},
			&shell.VoidFunction{
				Function: Vet,
			},
			&shell.BoolFunction{
				Function: Lint,
				Arg:      false,
			},
			&shell.VoidFunction{
				Function: Test,
			},
			&shell.VoidFunction{
				Function: Install,
			},
		},
	)
}

func Log() error {
	return shell.OrderedRunner(
		[]shell.RunningFunction{
			&shell.StringFunction{
				Arg:      "git log --all --decorate --oneline --graph",
				Function: shell.PrettyRun,
			},
		})
}
