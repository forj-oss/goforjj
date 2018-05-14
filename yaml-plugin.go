package goforjj

const ObjectApp = "app"

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
	YamlPluginTasksObjects `yaml:",inline"`
	instancesDetails map[string]*YamlPluginTasksObjects
}

func NewYamlPlugin() (ret *YamlPlugin) {
	ret = new(YamlPlugin)
	ret.instancesDetails = make(map[string]*YamlPluginTasksObjects)
	return 
}


func (*YamlPlugin) LoadIntance(from string) {
	
}
