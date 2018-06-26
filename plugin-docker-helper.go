package goforjj

import (
	"os"
	"strings"
)

func dockerCmd() (commands []string) {

	sudo := os.Getenv("DOCKER_SUDO")
	size := 1
	if sudo != "" {
		size++
	}
	commands = make([]string, 0, size)

	if sudo != "" {
		commands = append(commands, strings.Split(sudo, " ")...)
	}
	commands = append(commands, "docker")
	return
}
