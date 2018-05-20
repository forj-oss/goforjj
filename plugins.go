package goforjj

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// Plugins define a structure to store all plugins loaded.
type Plugins struct {
	drivers map[string]*Driver // key: plugin type, plugin name
	plugins map[string]map[string]*YamlPlugin
}

// NewPlugins create the list of plugins in memory.
func NewPlugins() (ret *Plugins) {
	ret = new(Plugins)
	ret.drivers = make(map[string]*Driver)
	ret.plugins = make(map[string]map[string]*YamlPlugin)
	return
}

// Load the instance plugin and return the driver object
// All plugin instances can share the same plugin. So that a plugin definition is loaded only once.
// But each driver are instance unique.
func (ps *Plugins) Load(instanceName, pluginName, pluginType string, loader map[string]func(*YamlPlugin) (yaml_data []byte, err error)) (driver *Driver, err error) {
	// check if driver already loaded
	if d, found := ps.drivers[instanceName]; found {
		return d, nil
	}

	plugin, newStatus := ps.definePlugin(pluginName, pluginType)
	var yaml_data []byte
	if newStatus {
		if lfunc, found := loader["master"] ; !found {
			delete(ps.plugins[pluginType], pluginName)
			return nil, fmt.Errorf("Internal issue. Load requires 'master' function loader")
		} else {
			yaml_data, err = lfunc(plugin)
		}
		
		if err != nil {
			delete(ps.plugins[pluginType], pluginName)
			plugin = nil
			return
		}
		if err = yaml.Unmarshal(yaml_data, plugin); err != nil {
			delete(ps.plugins[pluginType], pluginName)
			return
		}
	}

	driver = NewDriver(plugin)

	ps.drivers[instanceName] = driver

	// TODO: Need to load instances details now.
	if lfunc, found := loader["extended"] ; !found {
		return
	} else {
		yaml_data, err = lfunc(plugin)
	}
	
	if err != nil {
		plugin = nil
		return
	}

	if len(yaml_data) == 0 {
		return
	}
	pluginExtend := new(YamlPluginTasksObjects)
	if err = yaml.Unmarshal(yaml_data, pluginExtend); err != nil {
		return
	}

	driver.Yaml = driver.Yaml.MergeWith(instanceName, pluginExtend)

	return
}

func (ps *Plugins) definePlugin(driverName, driverType string) (plugin *YamlPlugin, new bool) {
	new = true
	if pn, ptFound := ps.plugins[driverType]; !ptFound {
		pt := make(map[string]*YamlPlugin)
		plugin = NewYamlPlugin()
		pt[driverName] = plugin
		ps.plugins[driverType] = pt
	} else if p, pFound := pn[driverName]; !pFound {
		plugin = NewYamlPlugin()
		pn[driverName] = plugin
		ps.plugins[driverType] = pn
	} else {
		plugin = p
		new = false
	}
	return
}
