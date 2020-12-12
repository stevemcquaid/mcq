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

	"github.com/stevemcquaid/mcq/pkg/colorWriter"
)

const ShellToUse = "bash"

func PrettyRun(command string) {
	yellowColorWriter := colorWriter.NewPrefixWriter(os.Stdout, color.New(color.FgCyan))
	defer yellowColorWriter.Flush()
	fmt.Fprintf(yellowColorWriter, "===> %s\n", command)

	greenColorWriter := colorWriter.NewPrefixWriter(os.Stdout, color.New(color.FgGreen))
	defer greenColorWriter.Flush()
	redColorWriter := colorWriter.NewPrefixWriter(os.Stdout, color.New(color.FgRed))
	defer redColorWriter.Flush()

	stdOutWriter := textio.NewPrefixWriter(greenColorWriter, "||    ")
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
}
