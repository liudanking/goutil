package osutil

import (
	"os"
	"os/exec"
)

type stdoutFilter struct{}

func (sf *stdoutFilter) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}

func RunCMD(cmdName string, args []string) ([]byte, error) {
	cmd := exec.Command(cmdName, args...)
	ret, err := cmd.CombinedOutput()
	return ret, err
}

func RunCMDFollow(cmdName string, args []string) error {
	cmd := exec.Command(cmdName, args...)
	sf := &stdoutFilter{}
	cmd.Stdout = sf
	cmd.Stderr = sf
	return cmd.Run()
}
