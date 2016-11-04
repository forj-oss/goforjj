package goforjj

import (
	"fmt"
	"net/http"
	"syscall"
)

func DefaultQuit(w http.ResponseWriter, ExitMessage string) {
	if ExitMessage == "" {
		ExitMessage = "Exiting"
	}
	fmt.Fprintln(w, ExitMessage)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

// Predefined Ping handler.
func PingHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "OK")
}
