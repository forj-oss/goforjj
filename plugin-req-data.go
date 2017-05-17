package goforjj

import (
	"encoding/json"
	"go/types"
)

//***************************************
// JSON data structure of plugin input.
// See plugin-actions.go about how those structs are managed.

type PluginReqData struct {
	// Collection of forjj flags requested by the plugin or given by default by Forjj
	Forj map[string]string
	// Define the list of Forjj objects data transmitted. object_type, instance, action.
	Objects map[string]ObjectInstances
}

type ObjectInstances map[string]InstanceKeys
type InstanceKeys map[string]ValueStruct
type ValueStruct struct {
	internal_type string
	value string
	list []string
} // Represents a flag value. Can be of type string or []string

func (v ValueStruct)MarshalJSON() ([]byte, error) {
	switch v.internal_type {
	case "string":
		return json.Marshal(v.value)
	case "[]string":
		return json.Marshal(v.list)
	}
	return nil
}

func (v *ValueStruct)Set(value interface{}) (ret *ValueStruct) {
	if v == nil {
		ret = new(ValueStruct)
	} else {
		ret = v
	}
	switch value.(type) {
	case string:
		ret.internal_type = "string"
		ret.value = value.(string)
	case []string:
		ret.internal_type = "[]string"
		ret.value = value.([]string)
	}
	return
}

func (v *ValueStruct)Get() (value interface{}) {
	switch v.internal_type {
	case "string":
		value = v.value
	case "[]string":
		value = v.list
	}
	return
}

func NewReqData() (r *PluginReqData) {
	r = new(PluginReqData)
	r.Forj = make(map[string]string)
	r.Objects = make(map[string]ObjectInstances)
	return
}

func (r *PluginReqData) SetForjFlag(key, value string) {
	if r == nil {
		return
	}
	if r.Forj == nil {
		r.Forj = make(map[string]string)
	}
	r.Forj[key] = value
}

func (r *PluginReqData) AddObjectActions(object_type, object_name string, keys InstanceKeys) {
	if r == nil {
		return
	}
	if r.Objects == nil {
		r.Objects = make(map[string]ObjectInstances)
	}
	if _, found := r.Objects[object_type]; !found {
		r.Objects[object_type] = make(map[string]InstanceKeys)
	}
	r.Objects[object_type][object_name] = keys
	return
}
