package commands

import (
	"fmt"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/shell"
)

func DockerBuild(dockerImage string) {
	shell.PrettyRun(fmt.Sprintf("docker build --target final -t %s .", dockerImage))
}

// @TODO - figure out port requirements
func DockerRun(dockerImage string) {
	DockerBuild(dockerImage)
	shell.PrettyRun(fmt.Sprintf("docker run -it -P %s .", dockerImage))
}

// @TODO - figure out port requirements
func DockerPush(dockerImage string) {
	DockerBuild(dockerImage)
	shell.PrettyRun(fmt.Sprintf("docker push %s", dockerImage))
}

// builds the binary
func Build(binaryName string) {
	shell.PrettyRun(fmt.Sprintf("go build -o bin/%s ./", binaryName))
	shell.PrettyRun(fmt.Sprintf("chmod u+x bin/%s", binaryName))
}

func BuildLinux(binaryName string) {
	shell.PrettyRun(fmt.Sprintf("GOOS=linux GOARCH=amd64 go build -o bin/%s ./", binaryName))
	shell.PrettyRun(fmt.Sprintf("chmod u+x bin/%s", binaryName))
}

func BuildWindows(binaryName string) {
	shell.PrettyRun(fmt.Sprintf("GOOS=windows GOARCH=386 go build -o bin/%s ./", binaryName))
	shell.PrettyRun(fmt.Sprintf("chmod u+x bin/%s", binaryName))
}

func Run() {
	shell.PrettyRun("go run main.go")
}

func Install() {
	shell.PrettyRun("go install")
}

func Vet() {
	shell.PrettyRun("go vet $(go list ./... | grep -v vendor) | grep -v '.pb.go:' | tee /dev/stderr")
}

func Clean() {
	Fmt()
	shell.PrettyRun("go mod tidy")
	Deps()
	Vet()
}

func Fmt() {
	shell.PrettyRun("find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s -l \"$file\"; done")
	shell.PrettyRun("find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w -l \"$file\"; done")
}

func Deps() {
	shell.PrettyRun("go mod tidy")
	shell.PrettyRun("go mod download")
	shell.PrettyRun("go mod vendor")
}

// Install all the build and lint dependencies
func Setup() {
	shell.PrettyRun("GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.32.2")
	shell.PrettyRun("GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck@2020.1.6")
	shell.PrettyRun("go get -u github.com/pierrre/gotestcover")
	shell.PrettyRun("go get -u golang.org/x/tools/cmd/cover")
	shell.PrettyRun("go get -u golang.org/x/tools/cmd/goimports")
	shell.PrettyRun("mkdir -p build")
}

func Test() {
	TestUnit()
}

func TestUnit() {
	shell.PrettyRun("go test -cover -covermode=atomic -coverprofile=build/unit.out $(go list ./... | grep -v /vendor/)  -run . -timeout=2m")
}

func TestIntegration() {
	shell.PrettyRun("go test -cover -covermode=atomic -coverprofile=build/integration.out -tags=integration $(go list ./... | grep -v /vendor/) -run $(TEST_PATTERN) -timeout=2m")
}

// Run all the tests and opens the coverage report
func Cover() {
	Test()
	shell.PrettyRun("gocovmerge build/unit.out > build/all.out")
	shell.PrettyRun("go tool cover -html=build/all.out")
}

func StaticCheck() {
	shell.PrettyRun("staticcheck -fail -tests -checks=\"all,-ST1000,-ST1021,-ST1020\" ./...")
}

func GolangCI() {
	command := []string{
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
	shell.PrettyRun(strings.Join(command, " "))
}

// Run all linters
func Lint() {
	GolangCI()
	StaticCheck()
}

// Run all the tests and code checks
func All() {
	Setup()
	CI()
	dockerImage, err := GetDockerImage()
	if err != nil {
		fmt.Printf("\n \n \n GOT ERROR: %s", err)
		return
	}
	DockerBuild(dockerImage)
}

// Run all the tests and code checks
func CI() {
	Deps()
	Fmt()
	Vet()
	Lint()
	Test()
}
