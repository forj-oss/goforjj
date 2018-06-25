package goforjj

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/forj-oss/forjj-modules/trace"
)

// Generic start command.
// If return code is 0 return the command output
// If return code is 1 return empty string
// else return the error string in err
func cmd_run(cmd_args []string) (string, error) {
	gotrace.Trace("RUNNING: '%s'", strings.Join(cmd_args, "' '"))
	cmd := exec.Command(cmd_args[0], cmd_args[1:]...)
	var ret []byte

	ret, err := cmd.CombinedOutput()
	if err == nil {
		gotrace.Trace("RET: %s", string(ret))
		return string(ret), nil
	}

	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	switch status.ExitStatus() {
	case 0:
		gotrace.Trace("RET: %s", string(ret))
		return string(ret), nil
	case 1:
		gotrace.Trace("RET: %s", string(ret))
		return "", nil
	default:
		return "", fmt.Errorf("ERROR: '%s' returns: %d. %s\n", cmd_args[0], status.ExitStatus(), cmd.ProcessState.String())
	}
}
