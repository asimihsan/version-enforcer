package command

import "os/exec"

// RunCommand runs the command and returns the output and error.
func RunCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
