package goforjj

import "gopkg.in/alecthomas/kingpin.v2"

// Defines default internal structure to enhance in the plugin.
type PluginTask struct {
  Cmd *kingpin.CmdClause
  Flags map[string]*string// Values for commands flags.
}

type ForjjPluginApp struct {
  App *kingpin.Application  // The kingpin Application structure for CLI flags management.
  IsInfra *bool             // True when creating infra repositories
  Tasks map[string]PluginTask
  Flags map[string]*string  // Values for global flags
}

//***************************************
// JSON data structure

type PluginRepo struct {
  Name string          // name of the repository
  Upstream string      // upstream url
  Files []string       // List of files managed by the plugin
}

 type PluginService struct {
  Url map[string]string
}

type PluginData struct {
  Repos map[string]PluginRepo   // List of repository data
  Services []PluginService      // web service url. ex: https://github.hpe.com
}

type PluginResult struct {
 Data PluginData
 State_code uint      // 200 OK
 Status string        // Status message
 Error_message string // Error message
}

//***************************************

type PluginDef struct {
 Name string
 Image string
}

// Following is created at create time or loaded from update/maintain
// File to define and store in the infra repository.
type PluginsDefinition struct {
 Plugins map[string]PluginDef // Ex: plugins["upstream"] = "github"
 Flow string                  // Ex: flow = "github-PR". This will connect all tools to provide a github PR flow Ready to start.
}

// Define the plugin yaml structure
// Yaml data structure

type YamlPluginRuntime struct {
    Service_type string
    Image string              `yaml:"docker_image"`
}

type YamlFlagsOptions struct {
    Help string
    Required bool
}

type YamlPluginDef struct {
    Help string
    Flags map[string]YamlFlagsOptions
}

type YamlPlugin struct {
    Name               string `yaml:"plugin"` // driver name              (yaml:plugin)
    Version            string                 // driver version           (yaml:version)
    Description        string                 // driver description       (yaml:description)
    Runtime YamlPluginRuntime                 // driver Runtime section   (yaml:runtime/)
    Flags map[string]YamlPluginDef          // flags definition         (yaml:flags/)
}

//***************************************

func init() {
}
