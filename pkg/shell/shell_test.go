package shell

import (
	"testing"
)

func TestPrettyRun(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy path",
			args: args{
				command: "ls",
			},
		}, {
			name: "run tests",
			args: args{
				command: "go test -cover -covermode=atomic -coverprofile=build/unit.out $(go list ./... | grep -v /vendor/)  -run . -timeout=2m",
			},
		}, {
			name: "merge results",
			args: args{
				command: "gocovmerge build/unit.out > build/all.out",
			},
		}, {
			name: "check cover",
			args: args{
				command: "go tool cover -html=build/all.out",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
