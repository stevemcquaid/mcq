package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fatih/color"
	// "github.com/fatih/color"
	"github.com/segmentio/textio"

	"github.com/stevemcquaid/mcq/pkg/colorwriter"
)

const ShellToUse = "sh"

// @TODO - create different pretty printers without the runner command. and use them inside the prettyrun()
func PrettyRun(command string) error {
	greenColorWriter := colorwriter.NewPrefixWriter(os.Stdout, color.New(color.FgGreen))
	defer greenColorWriter.Flush()
	_, _ = fmt.Fprintf(greenColorWriter, "===> %s\n", command)

	blueColorWriter := colorwriter.NewPrefixWriter(os.Stdout, color.New(color.FgCyan))
	defer blueColorWriter.Flush()
	redColorWriter := colorwriter.NewPrefixWriter(os.Stdout, color.New(color.FgRed))
	defer redColorWriter.Flush()

	stdOutWriter := textio.NewPrefixWriter(blueColorWriter, "||    ")
	defer stdOutWriter.Flush()

	stdErrWriter := textio.NewPrefixWriter(redColorWriter, "||    ")
	defer stdErrWriter.Flush()

	cmd := exec.Command(ShellToUse, "-c", command)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(stdOutWriter, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(stdErrWriter, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(redColorWriter, "------ cmd.Run() failed ------")
		fmt.Fprintln(stdErrWriter, err)

		return err

		// outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
		// if outStr != "" {
		// 	fmt.Println("------ stdout ---")
		// 	fmt.Println(outStr)
		// }
		// if errStr != "" {
		// 	fmt.Println("------ stderr ---")
		// 	fmt.Println(errStr)
		// }
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
