package goforjj

// Data stored about the driver
type PluginDef struct {
    Result *PluginResult // Json data structured returned.
    Yaml YamlPlugin      // Yaml data definition
}

// Following is created at create time or loaded from update/maintain
// File to define and store in the infra repository.
type PluginsDefinition struct {
    Plugins map[string]PluginDef // Ex: plugins["upstream"] = "github"
    Flow    string               // Ex: flow = "github-PR". This will connect all tools to provide a github PR flow Ready to start.
}

func init() {
}
