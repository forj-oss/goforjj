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
	Name                   string `yaml:"plugin"`
	Version                string `yaml:",omitempty"`
	Description            string `yaml:",omitempty"`
	CreatedFile            string `yaml:"created_flag_file"`
	ExtendRelPath          string `yaml:"extend_relative_path,omitempty"`
	Runtime                YamlPluginRuntime
	YamlPluginTasksObjects `yaml:",inline"`
	instancesDetails       map[string]*YamlPluginTasksObjects
}

func NewYamlPlugin() (ret *YamlPlugin) {
	ret = new(YamlPlugin)
	ret.instancesDetails = make(map[string]*YamlPluginTasksObjects)
	return
}

// MergeWith creates a new Yamlplugin object built from a merged between this plugin and the extension.
// except instance Details, source data won't be modified as there is no object pointer
func (p *YamlPlugin) MergeWith(instance string, extended *YamlPluginTasksObjects) (merged *YamlPlugin) {
	merged = new(YamlPlugin)
	*merged = *p
	p.instancesDetails[instance] = extended
	merged.instancesDetails = nil

	// Add extended task flags for inexistent in source.
	for taskType, taskList := range extended.Tasks {
		if srcList, found := p.Tasks[taskType]; found {
			for taskName, task := range taskList {
				if _, foundTask := srcList[taskName]; !foundTask {
					if srcList == nil {
						srcList = make(map[string]YamlFlag)
					}
					task.extentSource = true
					srcList[taskName] = task
					merged.Tasks[taskType] = srcList
				}
			}
		}
	}

	// Add extended object flags for inexistent in source.
	for objectType, object := range extended.Objects {
		if srcObject, found := p.Objects[objectType]; found {
			for flagName, flag := range object.Flags {
				if _, foundFlag := srcObject.Flags[flagName]; !foundFlag {
					if srcObject.Flags == nil {
						srcObject.Flags = make(map[string]YamlFlag)
					}
					flag.extentSource = true
					srcObject.Flags[flagName] = flag
					merged.Objects[objectType] = srcObject
				}
			}
			// Add extended object group of flags for inexistent in source.
			for groupName, group := range object.Groups {
				if srcGroup, foundGroup := srcObject.Groups[groupName]; !foundGroup {
					if srcObject.Groups == nil {
						srcObject.Groups = make(map[string]YamlObjectGroup)
					}
					for flag_name, flag := range group.Flags {
						flag.extentSource = true
						group.Flags[flag_name] = flag
					}
					srcObject.Groups[groupName] = group
					merged.Objects[objectType] = srcObject
				} else {
					for flagName, flag := range group.Flags {
						if _, foundFlag := srcGroup.Flags[flagName]; !foundFlag {
							if srcObject.Flags == nil {
								srcObject.Flags = make(map[string]YamlFlag)
							}
							flag.extentSource = true
							srcGroup.Flags[flagName] = flag
							srcObject.Groups[groupName] = srcGroup
							merged.Objects[objectType] = srcObject
						}
					}
				}
			}

		}
	}

	return
}
