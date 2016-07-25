package goforjj

import (
    "encoding/json"
    "fmt"
//    "github.hpe.com/christophe-larsonneur/goforjj/trace"
//    "os"
//    "os/exec"
//    "strings"
//    "syscall"
)

func PluginNew() *PluginResult {
    var p PluginResult
    p.Data.Repos = make(map[string]PluginRepo)
    p.Data.Services = make([]PluginService, 1)
    return &p
}

// Function to Start a driver as container in daemon mode.
// The container ID will be registered internally to kill/remove them in case of forjj panic.
func (p *PluginResult) PluginDockerRun(plugin_type, image, action string, docker_opts []string, opts []string) {

    // Check if container exists
    if _, err := docker_container_stat(plugin_type) ; err == nil {
        // need to remove it
        docker_container_remove(plugin_type)
    }

    docker_container_run(plugin_type, image, action, docker_opts, opts)
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
