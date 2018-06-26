package goforjj

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/forj-oss/forjj-modules/trace"
)

type commandRun struct {
	command     []string // Command to start
	args        []string // Arrays of args to provide to the command
	socket_path string   // Path to store the socket file
	socket_file string
	envs        map[string]string // Collection of environment variable to set
}

func (c *commandRun) Init(cmd ...string) error {
	if c == nil {
		return fmt.Errorf("cmdData object is nil")
	}

	if len(cmd) == 0 {
		return fmt.Errorf("Command cannot be empty")
	}
	c.command = cmd
	c.args = []string{}
	c.envs = make(map[string]string)
	return nil
}

func (c *commandRun) SetArgs(args []string) {
	if c == nil {
		return
	}
	c.args = args
}

func (c *commandRun) AddEnv(name, value string) {
	if c == nil {
		return
	}

	c.envs[name] = value
}

// runFlowCmd execute a command and transmit output to dedicated function to display it properly anywhere needed.
func (c *commandRun) runFlow(outFct func(string), errFct func(string)) (err error) {
	if c == nil {
		return fmt.Errorf("cmdData is nil")
	}

	if c.command == nil {
		return fmt.Errorf("Command cannot be empty")
	}
	command := c.command[0]
	args := make([]string, 0, len(c.args)+len(c.command[1:]))
	args = append(args, c.command[1:]...)
	args = append(args, c.args...)

	cmd := exec.Command(command, args...)
	outReader, _ := cmd.StdoutPipe()
	errReader, _ := cmd.StderrPipe()

	cmd.Env = make([]string, 0, len(c.envs)+len(os.Environ()))
	cmd.Env = append(cmd.Env, os.Environ()...)
	iCount := len(os.Environ()) - 1
	for key, value := range c.envs {
		cmd.Env[iCount] = key + "=" + value
		iCount++
	}
	gotrace.Trace("RUNNING: %s '%s'", command, strings.Join(args, "' '"))

	go func() {
		outScanner := bufio.NewScanner(outReader)
		for outScanner.Scan() {
			outFct(outScanner.Text())
		}
	}()

	go func() {
		outScanner := bufio.NewScanner(errReader)
		for outScanner.Scan() {
			errFct(outScanner.Text())
		}
	}()

	// Execute command
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("ERROR could not spawn command. %s", err.Error())
		return
	}

	gotrace.Trace("Command done")
	if status := cmd.ProcessState.Sys().(syscall.WaitStatus); status.ExitStatus() != 0 {
		err = fmt.Errorf("\n%s ERROR: Unable to get process status - %d: %s", c.command, status.ExitStatus(), cmd.ProcessState.String())
	}
	return

}
