package goforjj

//***************************************
// JSON data structure of plugin input.
// See plugin-actions.go about how those structs are managed.

// PluginReqData define the API data request to send to forjj plugins
type PluginReqData struct {
	// Collection of Forjj flags requested by the plugin or given by default by Forjj
	Forj       map[string]string
	ForjExtent map[string]string `json:",omitempty"` // Extended Forjj flags
	// Define the list of Forjj objects data transmitted. object_type, instance, action.
	Objects map[string]ObjectInstances
	Creds   map[string]string `json:",omitempty"` // Contains credentials requested by the plugin for a specific action.
}

// ObjectInstances is a collection of instanceKeys
type ObjectInstances map[string]InstanceKeys

// InstanceKeys is a collection of key/values or under key "extent" a collection of key/values (intanceExtentKeys)
type InstanceKeys map[string]interface{}

// InstanceExtentKeys is the collection of key/values which is stored as "extent" in InstanceKeys.
type InstanceExtentKeys map[string]*ValueStruct

// NewReqData return an empty API request structure.
func NewReqData() (r *PluginReqData) {
	r = new(PluginReqData)
	r.Forj = make(map[string]string)
	r.Objects = make(map[string]ObjectInstances)
	return
}

// SetForjFlag initialize forj part of the request with key/value or extent key/value.
func (r *PluginReqData) SetForjFlag(key, value string, cred, extent bool) {
	if r == nil {
		return
	}
	if cred {
		if r.Creds == nil {
			r.Creds = make(map[string]string)
		}
		r.Creds[key] = value
		if extent {
			return // For compatibility, creds data are still kept in normal structure So the function do not exit, except for extent (New way)
		}
	}
	if !extent {
		if r.Forj == nil {
			r.Forj = make(map[string]string)
		}
		r.Forj[key] = value
	} else {
		if r.ForjExtent == nil {
			r.ForjExtent = make(map[string]string)
		}
		r.ForjExtent[key] = value

	}
}

// AddObjectActions add in the request, the collection of keys/values or extent/keys/values for each objects/instances
func (r *PluginReqData) AddObjectActions(objectType, objectName string, keys InstanceKeys, extent InstanceExtentKeys, creds map[string]string) {
	if r == nil {
		return
	}
	if r.Objects == nil {
		r.Objects = make(map[string]ObjectInstances)
	}
	if _, found := r.Objects[objectType]; !found {
		r.Objects[objectType] = make(map[string]InstanceKeys)
	}
	keys["extent"] = extent
	r.Objects[objectType][objectName] = keys

	for cname, cdata := range creds {
		if r.Creds == nil {
			r.Creds = make(map[string]string)
		}
		r.Creds[objectType+"-"+objectName+"-"+cname] = cdata
	}
	return
}
