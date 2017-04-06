package goforjj

import "reflect"

const setup="setup"

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
