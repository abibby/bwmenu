package dmenu

import (
	"bytes"
	"os/exec"
	"strings"
)

func Open(options []string) (string, error) {
	cmd := exec.Command("dmenu")
	cmd.Stdin = bytes.NewBufferString(strings.Join(options, "\n"))
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
