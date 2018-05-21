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
		driverName    = "name"
		driverName2   = "name2"
		driverType    = "type"
		instanceName  = "instance Name"
		instanceName2 = "instance Name2"
		instanceName3 = "instance Name3"
		instanceName4 = "instance Name4"
		commonTask    = "common"
		flag1         = "flag1"
		flag2         = "flag2"
		flag3         = "flag3"
		group1        = "group1"
		group2        = "group2"
		obj1          = "obj1"
		obj2          = "obj2"
	)

	if plugins == nil {
		return
	}

	// --- Run the test ---
	// Check if we can load an instance without 'master' loader. We should not.
	plugin, err := plugins.Load(instanceName, driverName, driverType, map[string]func(*YamlPlugin) ([]byte, error){})

	// --- Start testing ---
	if plugin != nil {
		t.Error("Expected Plugins.Load to return nil object. Got one.")
		return
	} else if err == nil {
		t.Errorf("Expected Plugins.Load to return an error. Got none.")
	} else if v1, f1 := plugins.plugins[driverType]; !f1 {
		t.Error("Expected Plugins.Load to identify object type internally. Got none.")
	} else if _, f2 := v1[driverName]; f2 {
		t.Error("Expected Plugins.Load to identify object type but no driverName internally. Got one.")
	}

	// --- Run the test ---
	// Check if we can load an instance with a basic data. We could.
	plugin, err = plugins.Load(instanceName, driverName, driverType,
		map[string]func(*YamlPlugin) ([]byte, error){
			"master": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test"
				return []byte(data), nil
			},
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
	// Check if a loader error is reported if we already loaded the instance. We should not have any errors.
	plugin, err = plugins.Load(instanceName, driverName, driverType,
		map[string]func(*YamlPlugin) ([]byte, error){
			"master": func(_ *YamlPlugin) ([]byte, error) {
				return []byte{}, fmt.Errorf("test")
			},
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
	// Check if loading the instance2 with a loader error return an error. The driver entry should not exist.
	plugin, err = plugins.Load(instanceName2, driverName2, driverType,
		map[string]func(*YamlPlugin) ([]byte, error){
			"master": func(*YamlPlugin) ([]byte, error) {
				return []byte{}, fmt.Errorf("test")
			},
		})

	// --- Start testing ---
	if plugin != nil {
		t.Error("Expected Plugins.Load to return driver nil. Got an object.")
	} else if err == nil {
		t.Error("Expected Plugins.Load to return an error. No error returned.")
	}

	// --- Setting test context ---
	// --- Run the test ---
	// Check if loading an instance2 with some data task and object data are properly loaded .
	plugin, err = plugins.Load(instanceName2, driverName2, driverType,
		map[string]func(*YamlPlugin) ([]byte, error){
			"master": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test\n" +
					"task_flags:\n " +
					"  common:\n" +
					"    flag1:\n" +
					"      help: blabla\n" +
					"objects:\n" +
					"  obj1:\n" +
					"    flags:\n" +
					"      flag1:\n" +
					"        help: blublu\n" +
					"    groups:\n" +
					"      group2:\n" +
					"        flags:\n" +
					"          flag3:\n" +
					"            help: blibli\n"

				return []byte(data), nil
			},
		})

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected Plugins.Load to return NO error. got one: %s", err)
	} else if plugin == nil {
		t.Error("Expected Plugins.Load to return driver. Got nil.")
	} else if common, f1 := plugin.Yaml.Tasks[commonTask]; !f1 {
		t.Errorf("Expected Plugins.Load to have the task flags under '%s'. Not found", commonTask)
	} else if flag1Value, f2 := common[flag1]; !f2 {
		t.Errorf("Expected Plugins.Load to have the flag1 under '%s'. Found it", commonTask)
	} else if v1 := flag1Value.Help; v1 != "blabla" {
		t.Errorf("Expected Plugins.Load to have the flag1.help = 'blabla'. Got '%s'", v1)
	} else if objValue, f3 := plugin.Yaml.Objects[obj1]; !f3 {
		t.Errorf("Expected Plugins.Load to have the object '%s'. Not found", obj1)
	} else if obj1Flag1Value, f4 := objValue.Flags[flag1]; !f4 {
		t.Errorf("Expected Plugins.Load to have the object flag '%s/%s'. Not found", obj1, flag1)
	} else if v2 := obj1Flag1Value.Help; v2 != "blublu" {
		t.Errorf("Expected Plugins.Load to have the object '%s/%s'.help = 'blublu'. Got '%s'", obj1, flag1, v2)
	} else if obj1Group2Value, f5 := objValue.Groups[group2]; !f5 {
		t.Errorf("Expected Plugins.Load to have the object group '%s/%s'. Not found", obj1, group2)
	} else if obj1Group1Flag3Value, f6 := obj1Group2Value.Flags[flag3]; !f6 {
		t.Errorf("Expected Plugins.Load to have the object group '%s/%s/%s'. Not found", obj1, group2, flag3)
	} else if v3 := obj1Group1Flag3Value.Help; v3 != "blibli" {
		t.Errorf("Expected Plugins.Load to have the object '%s/%s/%s'.help = 'blibli'. Got '%s'", obj1, group2, flag3, v3)
	}

	// --- Setting test context ---
	// --- Run the test ---
	// Check if an already loaded instance&plugin won't be updated by another load from master or extended
	plugin, err = plugins.Load(instanceName2, driverName2, driverType,
		map[string]func(*YamlPlugin) ([]byte, error){
			"master": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test\n" +
					"task_flags:\n " +
					"  common:\n" +
					"    flag1:\n" +
					"      help: bloblo\n"
				return []byte(data), nil
			},
			"extended": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test\n" +
					"task_flags:\n " +
					"  common:\n" +
					"    flag2:\n" +
					"      help: blabla\n"
				return []byte(data), nil
			},
		})

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected Plugins.Load to return NO error. got one: %s", err)
	} else if plugin == nil {
		t.Error("Expected Plugins.Load to return driver. Got nil.")
	} else if common, f1 := plugin.Yaml.Tasks[commonTask]; !f1 {
		t.Errorf("Expected Plugins.Load to have the task flags under '%s'. Not found", commonTask)
	} else if flag1Value, f2 := common[flag1]; !f2 {
		t.Errorf("Expected Plugins.Load to load 'master'. But flag1 not found")
	} else if v2 := flag1Value.Help; v2 != "blabla" {
		t.Errorf("Expected Plugins.Load to load 'master' and 'extended' but flag1 created by 'master' (help:blabla) should not be updated by 'extended'. But flag1.help is '%s'", v2)
	} else if _, f3 := common[flag2]; f3 {
		t.Errorf("Expected Plugins.Load to load 'master' and 'extended' only. Loaded '%s'", "blabla")
	}

	// --- Setting test context ---
	// --- Run the test ---
	// check if instance3 with an already loaded plugin can load only extended. But others are ignored. New objects cannot be added by the extended. But new groups can.
	plugin, err = plugins.Load(instanceName3, driverName2, driverType,
		map[string]func(*YamlPlugin) ([]byte, error){
			"master": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test\n" +
					"task_flags:\n " +
					"  common:\n" +
					"    flag1:\n" +
					"      help: blibli"
				return []byte(data), nil
			},
			"blabla": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test\n" +
					"task_flags:\n " +
					"  common:\n" +
					"    flag2:\n" +
					"      help: blabla"
				return []byte(data), nil
			},
			"extended": func(*YamlPlugin) ([]byte, error) {
				data := "plugin: test\n" +
					"task_flags:\n " +
					"  common:\n" +
					"    flag1:\n" +
					"      help: bloblo\n" +
					"    flag2:\n" +
					"      help: bloblo\n" +
					"objects:\n" +
					"  obj1:\n" +
					"    flags:\n" +
					"      flag3:\n" +
					"        help: blibli\n" +
					"    groups:\n" +
					"      group1:\n" +
					"        flags:\n" +
					"          flag3:\n" +
					"            help: blibli\n" +
					"      group2:\n" +
					"        flags:\n" +
					"          flag1:\n" +
					"            help: blabla\n" +
					"  obj2:\n" +
					"    flags:\n" +
					"      flag3:\n" +
					"        help: blibli\n"
				return []byte(data), nil
			},
		})

	// --- Start testing ---
	if plugin == nil {
		t.Error("Expected Plugins.Load to return driver. Got nil.")
	} else if err != nil {
		t.Errorf("Expected Plugins.Load to return NO error. got one: %s", err)
	} else if common, f1 := plugin.Yaml.Tasks[commonTask]; !f1 {
		t.Errorf("Expected Plugins.Load to have the task flags under '%s'. Not found", commonTask)
	} else if flag1Value, f2 := common[flag1]; !f2 {
		t.Errorf("Expected Plugins.Load to load 'master'. But flag1 not found")
	} else if v2 := flag1Value.Help; v2 != "blabla" {
		t.Errorf("Expected Plugins.Load to have flag1.help already set to 'blabla', because driver is already loaded. But flag1.help is '%s'", v2)
	} else if flag2Value, f3 := common[flag2]; !f3 {
		t.Errorf("Expected Plugins.Load to load task flag2. Not found")
	} else if v3 := flag2Value.Help; v3 != "bloblo" {
		t.Errorf("Expected Plugins.Load to load 'extended' flag2 help 'blabla'. flag2.help is '%s'", v3)
	} else if _, f4 := plugin.Yaml.Objects[obj2]; f4 {
		t.Errorf("Expected Plugins.Load to not load undefined object 'obj2'. Got one")
	} else if obj1Value, f5 := plugin.Yaml.Objects[obj1]; !f5 {
		t.Errorf("Expected Plugins.Load to have object '%s'. Got none", obj1)
	} else if obj1Flag3, f6 := obj1Value.Flags[flag3]; !f6 {
		t.Errorf("Expected Plugins.Load to have object '%s-%s'. Got none", obj1, flag3)
	} else if v4 := obj1Flag3.Help; v4 != "blibli" {
		t.Errorf("Expected Plugins.Load to have object '%s-%s'.help = 'blibli'. Got '%s'", obj1, flag3, v4)
	} else if obj1Group1Value, f7 := obj1Value.Groups[group1]; !f7 {
		t.Errorf("Expected Plugins.Load to have object group '%s-%s'. Not found", obj1, group1)
	} else if obj1Group1Flag3Value, f8 := obj1Group1Value.Flags[flag3]; !f8 {
		t.Errorf("Expected Plugins.Load to have object group flag '%s-%s-%s'. Not found", obj1, group1, flag3)
	} else if v5 := obj1Group1Flag3Value.Help; v5 != "blibli" {
		t.Errorf("Expected Plugins.Load to have object group flag '%s-%s-%s'.help = 'blibli'. got '%s'", obj1, group1, flag3, v5)
	} else if obj1Group2Value, f9 := obj1Value.Groups[group2]; !f9 {
		t.Errorf("Expected Plugins.Load to have object group flag '%s-%s'. Not found", obj1, group2)
	} else if obj1Group2Flag1Value, f10 := obj1Group2Value.Flags[flag1]; !f10 {
		t.Errorf("Expected Plugins.Load to have object group flag '%s-%s-%s'. Not found", obj1, group2, flag1)
	} else if v6 := obj1Group2Flag1Value.Help; v6 != "blabla" {
		t.Errorf("Expected Plugins.Load to have object group flag '%s-%s-%s'.help = 'blabla'. got '%s'", obj1, group2, flag1, v6)
	}
}
