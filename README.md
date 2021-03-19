# mcq
A golang development helper. Dont memorize commands when you can `mcq lint`

# Usage
`mcq help`
```
$ mcq help
This application provides shortcuts to common development tasks

Usage:
  mcq [command]

Available Commands:
  all         Run everything
  build       -> go build
  ci          Run almost everything
  clean       -> fmt deps vet
  cover       -> go tool cover
  deps        -> go mod tidy, download, vendor
  docker      docker build, run, push
  fmt         -> go fmt
  help        Help about any command
  install     -> go install
  lint        -> golangci-lint, staticcheck
  log         -> ~git log --graph --oneline --decorate --all
  run         -> go run main.go
  setup       install dependencies
  test        -> go test
  version     Version

Flags:
  -h, --help   help for mcq

Use "mcq [command] --help" for more information about a command.
```

# TODO
* [x] Mechanism to fail fast during commands running. If error, it should quit. (OrderedRunner)
* [ ] Mechanism for pretty printing text to screen. Likely a writer library/passed around with global defaults for different types of messages
* [ ] Mechanism for parallelization of tasks than can be completed together
* [ ] Simplify colorwriter
