package goforjj

import "github.com/alecthomas/kingpin"

// Defines default internal structure to enhance in the plugin.
type PluginTask struct {
	Cmd   *kingpin.CmdClause
	Flags map[string]*string // Values for commands flags.
}

type ForjjPluginApp struct {
	App     *kingpin.Application // The kingpin Application structure for CLI flags management.
	IsInfra *bool                // True when creating infra repositories
	Tasks   map[string]PluginTask
	Flags   map[string]*string // Values for global flags
}
