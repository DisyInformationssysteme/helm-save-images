package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isDebug() bool {
	return os.Getenv("HELM_DEBUG") == "true"
}

func debugPrint(format string, a ...interface{}) {
	if isDebug() {
		fmt.Printf(format+"\n", a...)
	}
}

func outputWithRichError(cmd *exec.Cmd) ([]byte, error) {
	debugPrint("Executing %s", strings.Join(cmd.Args, " "))
	output, err := cmd.Output()
	if exitError, ok := err.(*exec.ExitError); ok {
		return output, fmt.Errorf("%s: %s", exitError.Error(), string(exitError.Stderr))
	}
	return output, err
}
