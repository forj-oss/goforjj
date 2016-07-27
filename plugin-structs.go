package goforjj

import (
    "github.com/parnurzeal/gorequest"
    "net/url"
)

// Data stored about the driver
type PluginDef struct {
    Result *PluginResult      // Json data structured returned.
    Yaml YamlPlugin           // Yaml data definition
    Source_path string        // Plugin source data
    service bool              // True if the service is started as daemon
    docker docker_container   // Define data to start the plugin as docker container
    cmd cmd_data              // Define data to start the service process
    req *gorequest.SuperAgent // REST API request
    url *url.URL              // REST API url
}

type cmd_data struct {
    command string       // Command to start
    args []string        // Arrays of args to provide to the command
    socket_path string   // Path to store the socket file
    socket_file string
}

type docker_container struct {
    name string
    opts []string
    socket_path string
}

// Following is created at create time or loaded from update/maintain
// File to define and store in the infra repository.
type PluginsDefinition struct {
    Plugins map[string]PluginDef // Ex: plugins["upstream"] = "github"
    Flow    string               // Ex: flow = "github-PR". This will connect all tools to provide a github PR flow Ready to start.
}

func init() {
}
