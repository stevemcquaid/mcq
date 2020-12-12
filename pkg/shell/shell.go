package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const ShellToUse = "bash"

func PrettyRun(command string){
	fmt.Printf("===> %s\n", command)

	cmd := exec.Command(ShellToUse, "-c", command)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		fmt.Println("------ cmd.Run() failed ---")
		fmt.Println("------ error --- ")
		fmt.Println(err)

		outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
		if outStr != "" {
			fmt.Println("------ stdout ---")
			fmt.Println(outStr)
		}
		if errStr != "" {
			fmt.Println("------ stderr ---")
			fmt.Println(errStr)
		}
	}
}
