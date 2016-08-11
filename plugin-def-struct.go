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
    Runtime     YamlPluginRuntime
    Actions     map[string]YamlPluginDef
}

// data structure in /flags
// actions: hash - Collection of valid keys. Support only common, create, update and maintain.
//   common: struct - Represent a collection of flags shared between create/update/maintain action
//     flags:
//   create: struct - Represent help and collection of flags for create action
//     help: string - Help attached to the action
//     flags: Hash - Collection of flags
//   update: struct - Represent help and collection of flags for update action
//     ... - same as create
//   maintain: struct - Represent help and collection of flags for maintain action
//     ... - same as create
type YamlPluginDef struct {
    Help  string
    Flags map[string]YamlFlagsOptions
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
    Image        string          `yaml:"docker_image"`
    Service      YamlPluginComm
    Shell        YamlPluginShell // Not yet used
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
    Socket     string
    Port       uint   // Not yet implemented
    Parameters []string // Not yet implemented
}

// data structure in /actions/<action>/flags/<flag name>
//     flags:
//       <flag name>:
//         help: string - Help attached to the flag
//         required: bool - true if this flag is required.
type YamlFlagsOptions struct {
    Help     string
    Required bool
    Hidden   bool // Used by the plugin.
    Default  string // Used by the plugin.
    Group    string // Group name used to regroup some flags under a dedicated struct.
                    // This group used only for create/update flags.
}
