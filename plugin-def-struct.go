package goforjj

// Yaml data structure

// Data structure in /
// ---
// plugin: string - Driver name (Name)
// version: string - driver version
// description: string - driver description
// runtime: struct - See YamlPluginRuntime
// actions: hash of struct - See YamlPluginDef - must be common/create/update/maintain as hash keys only.
type YamlPlugin struct {
	Name        string `yaml:"plugin"`
	Version     string
	Description string
	CreatedFile string `yaml:"created_flag_file"`
	Runtime     YamlPluginRuntime
	Tasks       map[string]map[string]YamlFlags `yaml:"task_flag"`
	Objects     map[string]YamlObject
}

// data structure in /objects/<Object Name>
//     flags:
//       <flag name>:
//         help: string - Help attached to the object
//         actions: collection of forjj actions (add/update/rename/remove/list)
type YamlObject struct {
	Actions            []string // Collection of actions for the group given.
	Help               string
	Identified_by_flag string // Multiple object configuration. each instance will have a key from a flag value
	Flags              map[string]YamlFlags
}

// data structure in /objects/<Object Name>/flags/<flag name>
//     flags:
//       <flag name>:
//         help: string - Help attached to the flag
//         required: bool - true if this flag is required.
type YamlFlags struct {
	Options YamlFlagOptions
	Help    string
	FormatRegexp string 	`yaml:"format-regexp"`
	Actions []string 		`yaml:"only-for-actions"`
}

type YamlFlagOptions struct {
	Required bool
	Hidden   bool   // Used by the plugin.
	Default  string // Used by the plugin.
	Secure   bool   // true if the data must be securely stored, ie not in the git repo. The flag must be defined in 'common' or 'maintain' flag group.
}

// data structure in /runtime
// runtime:
//   service_type: string - Support "REST API" and "shell"
//                          REST API means the driver comply to REST API served as web service
//                          shell means, the driver is called as shell with parameters and return a json data.
//   image_docker: string - Docker image containing the driver to start
//
type YamlPluginRuntime struct {
	Service_type string
	Docker       DockerStruct    `yaml:",omitempty"`
	Service      YamlPluginComm  `yaml:",omitempty"`
	Shell        YamlPluginShell `yaml:",omitempty"` // Not yet used
}

//data structure in /runtime/docker

type DockerStruct struct {
	Image   string
	Dood    bool              `yaml:",omitempty"`
	Volumes []string          `yaml:",omitempty"`
	Env     map[string]string `yaml:",omitempty"`
	User    string            `yaml:",omitempty"`
}

// data structure in /runtime/shell
// Define list of options for a shell plugin
// runtime:
//   shell:
//     parameters: Array of strings - List of parameters to provide to the shell/binary
//
type YamlPluginShell struct {
	Parameters []string
}

// data structure in /runtime/service
// 'service' defines how forjj communicate with the driver
// If service is not defined, socket will be used.
// runtime:
//   service:
//     socket: string - default set to the driver name with '.sock' as extension.
//                      The socket path is set by forjj.
//                      Socket file name to use between forjj and the driver
//     port: uint     - Port used to communicate between forjj and the driver
//     parameters: Array of strings - List of parameters to provide to the shell/binary
//                      Support {{Socket}}
type YamlPluginComm struct {
	Socket     string   `yaml:",omitempty"`
	Port       uint     `yaml:",omitempty"`      // Not yet implemented
	Parameters []string `yaml:",omitempty,flow"` // Not yet implemented
}
