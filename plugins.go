package goforjj

import (
	"gopkg.in/yaml.v2"
)

// Plugins define a structure to store all plugins loaded.
type Plugins struct {
	drivers map[string]*Plugin // key: plugin type, plugin name
	plugins map[string]map[string]*YamlPlugin
}

// NewPlugins create the list of plugins in memory.
func NewPlugins() (ret *Plugins) {
	ret = new(Plugins)
	ret.drivers = make(map[string]*Plugin)
	ret.plugins = make(map[string]map[string]*YamlPlugin)
	return
}

// Load the instance plugin and return the driver object
// All plugin instances can share the same plugin. So that a plugin definition is loaded only once.
// But each driver are instance unique.
func (ps *Plugins) Load(instanceName, driverName, driverType string, loader func() (yaml_data []byte, err error)) (driver *Plugin, err error) {
	// check if driver already loaded
	if d, found := ps.drivers[instanceName]; found {
		return d, nil
	}

	plugin, new := ps.definePlugin(driverName, driverType)
	if new {
		var yaml_data []byte
		yaml_data, err = loader()
		if err != nil {
			return
		}
		if err = yaml.Unmarshal(yaml_data, plugin); err != nil {
			return
		}
	}

	driver = NewDriver(plugin)

	ps.drivers[instanceName] = driver

	// TODO: Need to load instances details now.

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
