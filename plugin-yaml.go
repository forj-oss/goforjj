package goforjj

func NewYamlPlugin() (ret *YamlPlugin) {
	ret = new(YamlPlugin)
	ret.instancesDetails = make(map[string]*YamlPluginTasksObjects)
	return 
}

func (*YamlPlugin) LoadIntance() {
	
}