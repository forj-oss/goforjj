package goforjj

import (
	"fmt"
	"testing"
)

func TestNewPlugins(t *testing.T) {
	t.Log("Expect NewPlugins to initialize Plugins object.")

	// --- Setting test context ---
	// --- Run the test ---
	ret := NewPlugins()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected Plugins to return valid object. Got nil.")
		return
	}
	if ret.drivers == nil {
		t.Error("Expected Plugins to initialized properly. drivers is nil.")
		return
	} else if len(ret.drivers) != 0 {
		t.Error("Expected Plugins.drivers to be empty. is not.")
	}
	if ret.plugins == nil {
		t.Error("Expected Plugins to initialized properly. plugins is nil.")
		return
	} else if len(ret.plugins) != 0 {
		t.Error("Expected Plugins.plugins to be empty. is not.")
	}
}

func TestDefinePlugin(t *testing.T) {
	t.Log("Expect definePlugin to define properly object type and name.")

	// --- Setting test context ---
	ret := NewPlugins()

	const (
		driverName = "name"
		driverType = "type"
	)

	if ret == nil {
		return
	}

	// --- Run the test ---
	plugin, new_value := ret.definePlugin(driverName, driverType)

	// --- Start testing ---
	if plugin == nil {
		t.Error("Expected Plugins.definePlugin to return valid object. Got nil.")
		return
	}
	if plugin.Name != "" {
		t.Error("Expected plugin.Name to be uninitialized. Is not.")
	} else if len(ret.drivers) != 0 {
		t.Error("Expected Plugins.drivers to be empty. is not.")
	}

	if ret.plugins == nil {
		t.Error("Expected Plugins to initialized properly. plugins is nil.")
		return
	} else if v := len(ret.plugins); v != 1 {
		t.Errorf("Expected Plugins.plugins to have a single record. Got %d record(s).", v)
	} else if v1, found := ret.plugins[driverType]; !found {
		t.Errorf("Expected Plugins.plugins[%s] to exist. not found.", driverType)
	} else if v2, found2 := v1[driverName]; !found2 {
		t.Errorf("Expected Plugins.plugins[%s][%s] to exist. not found.", driverType, driverName)
	} else {
		v2.Name = "test"
		if plugin.Name != "test" {
			t.Error("Expected plugin to be stored in Plugins.plugins. Not the same object.")
		}
		if !new_value {
			t.Error("Expected plugin to be new. new is False")
		}
	}

	// --- Setting test context ---
	const (
		driverName2 = "name2"
	)

	// --- Run the test ---
	plugin, new_value = ret.definePlugin(driverName2, driverType)

	// --- Start testing ---
	if plugin == nil {
		t.Error("Expected Plugins.definePlugin to return valid object. Got nil.")
		return
	}
	if plugin.Name != "" {
		t.Error("Expected plugin.Name to be uninitialized. Is not.")
	} else if len(ret.drivers) != 0 {
		t.Error("Expected Plugins.drivers to be empty. is not.")
	}

	if ret.plugins == nil {
		t.Error("Expected Plugins to initialized properly. plugins is nil.")
		return
	} else if v := len(ret.plugins); v != 1 {
		t.Errorf("Expected Plugins.plugins to have a single record. Got %d record(s).", v)
	} else if v1, found := ret.plugins[driverType]; !found {
		t.Errorf("Expected Plugins.plugins[%s] to exist. not found.", driverType)
	} else if n := len(v1); n != 2 {
		t.Errorf("Expected Plugins.plugins[%s] to have 2 records. Got %d record(s).", driverType, n)
	} else if v2, found2 := v1[driverName2]; !found2 {
		t.Errorf("Expected Plugins.plugins[%s][%s] to exist. not found.", driverType, driverName)
	} else {
		v2.Name = "test2"
		if plugin.Name != "test2" {
			t.Error("Expected plugin to be stored in Plugins.plugins. Not the same object.")
		}
		if !new_value {
			t.Error("Expected plugin to be new. new is False")
		}
	}

	// --- Run the test ---
	plugin, new_value = ret.definePlugin(driverName, driverType)

	// --- Start testing ---
	if plugin == nil {
		t.Error("Expected Plugins.definePlugin to return valid object. Got nil.")
		return
	}
	if plugin.Name != "test" {
		t.Errorf("Expected plugin.Name to be already initialized properly. Got '%s'.", plugin.Name)
	} else if len(ret.drivers) != 0 {
		t.Error("Expected Plugins.drivers to be empty. is not.")
	}

	if ret.plugins == nil {
		t.Error("Expected Plugins to initialized properly. plugins is nil.")
		return
	} else if v := len(ret.plugins); v != 1 {
		t.Errorf("Expected Plugins.plugins to have a single record. Got %d record(s).", v)
	} else if v1, found := ret.plugins[driverType]; !found {
		t.Errorf("Expected Plugins.plugins[%s] to exist. not found.", driverType)
	} else if n := len(v1); n != 2 {
		t.Errorf("Expected Plugins.plugins[%s] to have 2 records. Got %d record(s).", driverType, n)
	} else if v2, found2 := v1[driverName]; !found2 {
		t.Errorf("Expected Plugins.plugins[%s][%s] to exist. not found.", driverType, driverName)
	} else {
		v2.Name = "test3"
		if plugin.Name != "test3" {
			t.Error("Expected plugin to be stored in Plugins.plugins. Not the same object.")
		}
		if new_value {
			t.Error("Expected plugin to be identified as existing. new is True")
		}
	}

}

func TestLoad(t *testing.T) {
	t.Log("Expect Load to load properly the plugin and return a driver object.")

	// --- Setting test context ---
	plugins := NewPlugins()

	const (
		driverName   = "name"
		driverName2  = "name2"
		driverType   = "type"
		instanceName = "instance Name"
		instanceName2 = "instance Name2"
	)

	if plugins == nil {
		return
	}

	// --- Run the test ---
	plugin, err := plugins.Load(instanceName, driverName, driverType, func() ([]byte, error) {
		data := "plugin: test"
		return []byte(data), nil
	})

	// --- Start testing ---
	if plugin == nil {
		t.Error("Expected Plugins.Load to return a valid plugin object. Got nil.")
		return
	} else if err != nil {
		t.Errorf("Expected Plugins.Load to return no error. Got '%s'.", err)
	} else if v1 := plugin.Yaml.Name; v1 != "test" {
		t.Errorf("Expected Plugins.Load to return a valid plugin name. Got '%s'.", v1)
	} else {
		v2 := plugins.plugins[driverType]
		v3 := v2[driverName]
		v3.Name = "test2"
		if plugin.Yaml.Name != "test2" {
			t.Error("Expected Plugins.Load to return a valid plugin name attached to a defined plugin. not the same plugin.")
		}
	}

	// --- Setting test context ---
	// --- Run the test ---
	plugin, err = plugins.Load(instanceName, driverName, driverType, func() ([]byte, error) {
		return []byte{}, fmt.Errorf("test")
	})

	// --- Start testing ---
	if plugin == nil {
		t.Error("Expected Plugins.Load to return the driver already loaded. Got nil.")
	} else if err != nil {
		t.Errorf("Expected Plugins.Load to return without failure. Got error '%s'.", err)
	} else {
		if plugin.Yaml.Name != "test2" {
			t.Error("Expected Plugins.Load to return the correct valid plugin/driver. Got a different one.")
		}
	}

	// --- Setting test context ---
	// --- Run the test ---
	plugin, err = plugins.Load(instanceName2, driverName2, driverType, func() ([]byte, error) {
		return []byte{}, fmt.Errorf("test")
	})

	// --- Start testing ---
	if plugin != nil {
		t.Error("Expected Plugins.Load to return driver nil. Got an object.")
	} else if err == nil {
		t.Error("Expected Plugins.Load to return an error. No error returned.")
	}

}
