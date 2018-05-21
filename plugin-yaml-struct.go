package goforjj

// Yaml data structure
// See root one in yaml-plugin.go


type YamlPluginTasksObjects struct {
	Tasks       map[string]map[string]YamlFlag `yaml:"task_flags"`
	Objects     map[string]YamlObject
}

// data structure in /objects/<Object Name>/lists/<list_name>
type YamlObjectList struct {
	Sep       string `yaml:"separator"`
	Help      string
	ExtRegexp string `yaml:"defined-by"`
}

// data structure in /objects/<Object Name>/groups/<group_name>
type YamlObjectGroup struct {
	Actions []string `yaml:"default-actions"` // Collection of actions for the group given.
	Flags   map[string]YamlFlag
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
	Docker       DockerStruct   `yaml:",omitempty"`
	Service      YamlPluginComm `yaml:",omitempty"`
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
	Command    string   `yaml:",omitempty"`      // Not yet implemented
	Parameters []string `yaml:",omitempty,flow"` // Not yet implemented
}

