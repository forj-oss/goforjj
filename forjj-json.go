package goforjj

import (
        "encoding/json"
        "fmt"
        "os"
        "os/exec"
        "strings"
        "syscall"
)

// Function to read json
func (p *PluginResult)PluginRun(image, action string, docker_opts []string, opts []string) {
 cmd_args := append([]string{}, "sudo", "docker", "run", "-i", "--rm")
 cmd_args = append(cmd_args, docker_opts ...)
 cmd_args = append(cmd_args, image, action)
 cmd_args = append(cmd_args, opts ...)

 cmd := exec.Command(cmd_args[0], cmd_args[1:]...)
 fmt.Printf("RUNNING: %s\n\n", strings.Join(cmd_args, " "))

 stdout, err := cmd.StdoutPipe()
 if err != nil {
   fmt.Println(err)
   os.Exit(1)
 }

 cmd.Start()

 if err := json.NewDecoder(stdout).Decode(p); err != nil {
   fmt.Println(err)
   os.Exit(1)
 }
 if err := cmd.Wait(); err != nil {
    fmt.Printf("\n%s ERROR.\nCommand status: %s\n", action, err)
    os.Exit(1)
 }

 if status := cmd.ProcessState.Sys().(syscall.WaitStatus) ; status.ExitStatus() != 0 {
    fmt.Printf("\n%s ERROR.\nCommand status: %s\n", action, cmd.ProcessState.String())
    os.Exit(status.ExitStatus())
 }
 println(action, "DONE")
}

// Function to print out json data
func (p *PluginResult)JsonPrint() error {
 if b, err := json.Marshal(p) ; err != nil {
   return err
 } else {
   fmt.Printf("%s\n", b)
 }
 return nil
}
