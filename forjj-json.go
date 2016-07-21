package goforjj

import (
    "encoding/json"
    "fmt"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
    "os"
    "os/exec"
    "strings"
    "syscall"
)

func PluginNew() *PluginResult {
    var p PluginResult
    p.Data.Repos = make(map[string]PluginRepo)
    p.Data.Services = make([]PluginService, 1)
    return &p
}

// Function to read json
func (p *PluginResult) PluginRun(image, action string, docker_opts []string, opts []string) {
    if image == "" {
        fmt.Printf("docker_image is missing in the driver definition. driver ignored.\n")
        return
    }
    cmd_args := append([]string{}, "sudo", "docker", "run", "-i", "--rm")
    cmd_args = append(cmd_args, docker_opts...)
    cmd_args = append(cmd_args, image, action)
    cmd_args = append(cmd_args, opts...)

    cmd := exec.Command(cmd_args[0], cmd_args[1:]...)
    gotrace.Trace("RUNNING: %s\n\n", strings.Join(cmd_args, " "))

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

    if status := cmd.ProcessState.Sys().(syscall.WaitStatus); status.ExitStatus() != 0 {
        fmt.Printf("\n%s ERROR.\nCommand status: %s\n", action, cmd.ProcessState.String())
        os.Exit(status.ExitStatus())
    }
    gotrace.Trace("%s %s\n", action, "DONE")
}

// Function to print out json data
func (p *PluginResult) JsonPrint() error {
    if b, err := json.Marshal(p); err != nil {
        return err
    } else {
        fmt.Printf("%s\n", b)
    }
    return nil
}

/*
// Load data returned by the plugin in the internal structure of Forjj core.
+func (p *PluginsDefinition)LoadResult(res *PluginResult) error {
+ if p.plugins == nil { p.plugins = make(map[string]PluginDef) }
+ return nil
 }
*/
