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

// JSON data structure

type PluginRepo struct {
  Name string          // name of the repository
  Upstream string      // upstream url
  Files []string       // List of files managed by the plugin
}

type PluginData struct {
  Repos map[string]PluginRepo   // List of repository data
  Services map[string]string    // web service url. ex: https://github.hpe.com
}

type PluginResult struct {
 data PluginData
 state_code uint      // 200 OK
 status string        // Status message
 error_message string // Error message
}

func init() {
}
