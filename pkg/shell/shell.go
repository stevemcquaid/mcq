package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	// "github.com/fatih/color"
	"github.com/segmentio/textio"

	"github.com/stevemcquaid/mcq/pkg/colorwriter"
)

const ShellToUse = "sh"

// createFlushDefer creates a defer function that safely flushes a writer
func createFlushDefer(writer io.Writer, writerName string) func() {
	return func() {
		if flushable, ok := writer.(interface{ Flush() error }); ok {
			if err := flushable.Flush(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error flushing %s: %v\n", writerName, err)
			}
		}
	}
}

// setupWriters creates and configures all the color writers needed for PrettyRun
func setupWriters() (greenWriter, blueWriter, redWriter, stdOutWriter, stdErrWriter io.Writer) {
	greenWriter = colorwriter.NewPrefixWriter(os.Stdout, color.New(color.FgGreen))
	blueWriter = colorwriter.NewPrefixWriter(os.Stdout, color.New(color.FgCyan))
	redWriter = colorwriter.NewPrefixWriter(os.Stdout, color.New(color.FgRed))

	stdOutWriter = textio.NewPrefixWriter(blueWriter, "||    ")
	stdErrWriter = textio.NewPrefixWriter(redWriter, "||    ")

	return greenWriter, blueWriter, redWriter, stdOutWriter, stdErrWriter
}

// executeCommand runs the command and handles the output
func executeCommand(command string, stdOutWriter, stdErrWriter io.Writer) error {
	cmd := exec.Command(ShellToUse, "-c", command)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(stdOutWriter, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(stdErrWriter, &stderrBuf)

	return cmd.Run()
}

// handleCommandError handles the case when command execution fails
func handleCommandError(err error, redWriter, stdErrWriter io.Writer) error {
	_, _ = fmt.Fprintln(redWriter, "------ cmd.Run() failed ------")
	_, _ = fmt.Fprintln(stdErrWriter, err)
	return err
}

// @TODO - create different pretty printers without the runner command. and use them inside the prettyrun()
func PrettyRun(command string) error {
	greenWriter, blueWriter, redWriter, stdOutWriter, stdErrWriter := setupWriters()

	// Set up defer functions for flushing
	defer createFlushDefer(greenWriter, "green color writer")()
	defer createFlushDefer(blueWriter, "blue color writer")()
	defer createFlushDefer(redWriter, "red color writer")()
	defer createFlushDefer(stdOutWriter, "stdout writer")()
	defer createFlushDefer(stdErrWriter, "stderr writer")()

	// Print the command being executed
	_, _ = fmt.Fprintf(greenWriter, "===> %s\n", command)

	// Execute the command
	err := executeCommand(command, stdOutWriter, stdErrWriter)
	if err != nil {
		return handleCommandError(err, redWriter, stdErrWriter)
	}

	return nil
}

// RunningFunction defines a generic interface to run functions
type RunningFunction interface {
	Run() error
}

// OrderedRunner takes an array of objects of type RunningFunction and tells each to run in sequence, quitting if there are any errors
func OrderedRunner(queue []RunningFunction) error {
	for _, item := range queue {
		err := item.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// OrderedRunnerIgnoreErrors takes an array of objects of type RunningFunction and tells each to run in sequence, it aggregates errors and provides all at the end
func OrderedRunnerIgnoreErrors(queue []RunningFunction) error {
	var myErrors []string
	for _, item := range queue {
		err := item.Run()
		if err != nil {
			myErrors = append(myErrors, err.Error())
		}
	}
	if myErrors != nil {
		return fmt.Errorf("error: %s", strings.Join(myErrors, ". Error: "))
	}
	return nil
}

// StringSliceFunction implements RunningFunction interface, and supports Functions with a single string argument
type StringSliceFunction struct {
	Arg      []string
	Function func(input []string) error
}

func (f *StringSliceFunction) Run() error {
	return f.Function(f.Arg)
}

// StringFunction implements RunningFunction interface, and supports Functions with a single string argument
type StringFunction struct {
	Arg      string
	Function func(input string) error
}

func (f *StringFunction) Run() error {
	return f.Function(f.Arg)
}

// VoidFunction implements RunningFunction interface, and supports Functions with no arguments
type VoidFunction struct {
	Function func() error
}

func (f *VoidFunction) Run() error {
	return f.Function()
}

// BoolFunction implements RunningFunction interface, and supports Functions with a single bool argument
type BoolFunction struct {
	Arg      bool
	Function func(input bool) error
}

func (f *BoolFunction) Run() error {
	return f.Function(f.Arg)
}
