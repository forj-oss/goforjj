package goforjj

import (

)

type cmd_data struct {
	command     string   // Command to start
	args        []string // Arrays of args to provide to the command
	socket_path string   // Path to store the socket file
	socket_file string
}

type docker_container struct {
	name        string
	opts        []string
	socket_path string
	volumes     map[string]byte
	envs        map[string]byte
}

// Following is created at create time or loaded from update/maintain
// File to define and store in the infra repository.
type PluginsDefinition struct {
	Plugins map[string]*Plugin   // Ex: plugins["upstream"] = "github"
	Flow    string               // Ex: flow = "github-PR". This will connect all tools to provide a github PR flow Ready to start.
}

func init() {
}
