package utils

import (
	"bytes"
	"os/exec"
	"runtime"
)

func ExecuteHostCmd(cmd string) (string, error) {

	var stdout, stderr bytes.Buffer
	var command *exec.Cmd
	if runtime.GOOS == "windows" { // Windows
		command = exec.Command("cmd", "/C", cmd)
	} else { // Linux/Unix
		command = exec.Command("sh", "-c", cmd)
	}
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return stderr.String(), err
	}

	return stdout.String(), nil
}
