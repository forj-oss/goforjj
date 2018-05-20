package goforjj

import (
	"reflect"
	"strings"
)

const setup="setup"

// YamlObject data structure in /objects/<Object Name>
//     flags:
//       <flag name>:
//         help: string - Help attached to the object
//         actions: collection of forjj actions (add/update/rename/remove/list)
type YamlObject struct { // Ex: object projects. Instances: prj1, prj2, ...
	Actions            []string `yaml:"default-actions"` // Collection of actions for the group given.
	Help               string
	FlagsScope         string `yaml:"flags-scope"` // 'object' by default. flag name is NOT prefixed
	// 'instance' flag name is prefixed by instance name.
	FieldsScope        string `yaml:"fields-scope"` // 'global' by default. Means field is added at Object level.
	// 'instance' Means fields is added at object instance level.
	Identified_by_flag string // Multiple object configuration. each instance will have a key from a flag value
	Groups             map[string]YamlObjectGroup
	Flags              map[string]YamlFlag
	Lists              map[string]YamlObjectList
}

func (o *YamlObject)found_flag_action_def(action, group_name, flag_name string) (found bool) {
	if o == nil {
		return
	}

	var actions []string
	var flag *YamlFlag

	if group_name == "" {
		if o.Flags == nil {
			return
		}
		if v, found := o.Flags[flag_name] ; !found {
			return false
		} else {
			flag = &v
		}
	} else {
		if o.Groups == nil {
			return
		}
		if v, found := o.Groups[group_name] ; !found {
			return false
		} else {
			if v.Flags == nil {
				return false
			}
			if f, found := v.Flags[flag_name] ; found {
				flag = &f
			} else {
				return false
			}
		}
	}
	if action == setup {
		return true
	}
	actions = flag.Actions
	if actions != nil && len(actions) > 0 {
		return flag.found_action_def(action)
	}
	if found, _ = InArray(action, actions) ; found {
		return
	}
	found, _ = InArray(setup, actions)
	return
}

func InArray(v interface{}, in interface{}) (ok bool, i int) {
    val := reflect.Indirect(reflect.ValueOf(in))
    switch val.Kind() {
    case reflect.Slice, reflect.Array:
        for ; i < val.Len(); i++ {
            if ok = v == val.Index(i).Interface(); ok {
                return
            }
        }
    }
    return
}

func (o *YamlObject)FlagsRange(action string) (res map[string]YamlFlag) {
	if o == nil {
		return
	}
	res = make(map[string]YamlFlag)
	for group_name, group := range o.Groups {
		for flag_name, flag := range group.Flags {
			if o.found_flag_action_def(action, group_name, flag_name) {
				res[group_name + "-" + flag_name] = flag
			}
		}
	}
	for flag_name, flag := range o.Flags {

		if o.found_flag_action_def(action, "", flag_name) {
			res[flag_name] = flag
		}
	}
	return
}

func (o *YamlObject)HasValidKey(key string) bool {
	if _, found := o.Flags[key]; found {
		return true
	}
	for group_name, group := range o.Groups {
		if len(group_name) + 1 < len(key) && strings.HasPrefix(key, group_name) {
			key_name := key[len(group_name) + 1:]
			if _, found := group.Flags[key_name] ; found {
				return true
			}
		}
	}
	return false
}
