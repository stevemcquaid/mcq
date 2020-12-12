package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

const ShellToUse = "bash"

func Run(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func PrettyRun(command string) {
	fmt.Printf("---> %s\n", command)
	out, errout, err := Run(command)
	if err != nil {
		fmt.Println("------ error --- ")
		fmt.Println(err)
	}
	if out != "" {
		fmt.Println("------ stdout ---")
		fmt.Println(out)
	}
	if errout != "" {
		fmt.Println("------ stderr ---")
		fmt.Println(errout)
	}
}
