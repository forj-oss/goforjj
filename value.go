package goforjj

import (
	"encoding/json"
	"fmt"
	"reflect"
	"text/template"
	"bytes"
	"strings"
	"github.com/forj-oss/forjj-modules/trace"
)

type ValueStruct struct {
	internal_type string
	value string
	list []string
} // Represents a flag value. Can be of type string or []string

func (v ValueStruct)MarshalJSON() ([]byte, error) {
	switch v.internal_type {
	case "[]string":
		return json.Marshal(v.list)
	}
	// By default, encode a string
	return json.Marshal(v.value)
}

func (v ValueStruct)MarshalYAML() (interface{}, error) {
	switch v.internal_type {
	case "[]string":
		return v.list, nil
	}
	// By default, encode a string
	return v.value, nil
}

func (v *ValueStruct) UnmarshalYAML(unmarchal func(interface{}) error) error {
	var data interface{}

	if err := unmarchal(&data); err != nil {
		return err
	}
	switch data.(type) {
	case string:
		v.internal_type = "string"
		v.value = data.(string)
	case []string:
		v.internal_type = "[]string"
		v.list = data.([]string)
	default:
		return fmt.Errorf("value type ('%s') not supported.", reflect.TypeOf(data))
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
		ret.SetString(value.(string))
	case []string:
		ret.internal_type = "[]string"
		ret.list = value.([]string)
	case *ValueStruct:
		*ret = *value.(*ValueStruct)
	case ValueStruct:
		*ret = value.(ValueStruct)
	}
	return
}

func (v *ValueStruct)Evaluate(data interface{}) error {
	var doc bytes.Buffer
	tmpl := template.New("forjj_data")

	switch v.internal_type {
	case "string":
		if ! strings.Contains(v.value, "{{") { return nil }
		if _, err := tmpl.Funcs(template.FuncMap{
			"ToLower" : strings.ToLower,
		}).Parse(v.value) ; err != nil {
			return err
		}
		if err := tmpl.Execute(&doc, data) ; err != nil {
			return err
		}
		ret := doc.String()
		gotrace.Trace("'%s' were interpreted to '%s'", v.value, ret)
		v.value = ret
	case "[]string":
		for index, value := range v.list {
			if ! strings.Contains(v.value, "{{") {
				continue
			}
			if _, err := tmpl.Funcs(template.FuncMap{
				"ToLower" : strings.ToLower,
			}).Parse(value); err != nil {
				return err
			}
			if err := tmpl.Execute(&doc, data) ; err != nil {
				return err
			}
			ret := doc.String()
			gotrace.Trace("'%s'[%d] were interpreted to '%s'", v.value, index, ret)
			v.list[index] = ret
		}
	}
	return nil
}

func (v *ValueStruct)SetIfFound(value interface{}, found bool) (ret *ValueStruct, ret_bool bool) {
	if !found {
		return
	}
	ret_bool = true
	ret = v.Set(value)
	return
}

func (v *ValueStruct)SetString(value string) {
	if v == nil { return }
	v.internal_type = "string"
	v.value = value

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

func (v *ValueStruct)Type() string {
	return v.internal_type
}

func (v *ValueStruct)GetString() (string) {
	if v == nil { return "" }
	if v.internal_type != "string" { return "" }
	return v.value
}

func (v *ValueStruct)GetStringSlice() ([]string) {
	if v == nil { return nil }
	if v.internal_type != "[]string" { return nil }
	return v.list
}

// Equal return true, if value compared are equal.
// The equality depends on internal type:
// - string : string equality
// - []string : same list of elements and each elements are at the same position
func (v *ValueStruct)Equal(value *ValueStruct) (bool) {
	if v == nil && value == nil {
		return true
	}
	if v == nil || value == nil {
		return false
	}
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
