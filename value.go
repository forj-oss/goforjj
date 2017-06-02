package goforjj

import (
	"encoding/json"
)

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
	return nil, nil
}

func (v *ValueStruct)Set(value interface{}, found bool) (ret *ValueStruct, ret_bool bool) {
	if !found {
		return
	}
	ret_bool = true
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
		ret.list = value.([]string)
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

// Equal return true, if value compared are equal.
// The equality depends on internal type:
// - string : string equality
// - []string : same list of elements and each elements are at the same position
func (v *ValueStruct)Equal(value *ValueStruct) (bool) {
	if v.internal_type != value.internal_type {
		return false
	}
	switch v.internal_type {
	case "string":
		return (v.value == value.value)
	case "[]string":
		if v.list == nil || value.list == nil { return (v.list == nil && value.list == nil) }
		if len(v.list) != len(value.list)     { return false }
		for index, element := range v.list {
			if value.list[index] != element   { return false }
		}
		return true
	default:
		return false
	}
}
