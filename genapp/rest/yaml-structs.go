package main

// ************************
// Create request structure
// ************************

type CreateReq struct {
	Forj struct {
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.common }}\
		Instance string `json:"instance"` // __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		// __MYPLUGIN: {{ end }}\
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.create }}\
		// __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		// __MYPLUGIN: {{ end }}\
	}
	Objects CreateArgReq
}

type CreateArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $ObjectOpts := .Yaml.Objects }}\
	App map[string]AppInstanceStruct // __MYPLUGIN: 	{{ go_vars $Objectname}} map[string]{{ go_vars $Objectname}}InstanceStruct `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{ end }}\
}

// __MYPLUGIN: {{ range $Objectname, $ObjectOpts := .Yaml.Objects }}\
type AppInstanceStruct struct {
	// __MYPLUGIN: type {{ go_vars $Objectname}}InstanceStruct struct {
	// __MYPLUGIN: {{   range $ActionName, $Flags := object_tree $ObjectOpts }}\
	Action AppActionStruct // __MYPLUGIN: 	{{ go_vars $ActionName}} {{ go_vars $Objectname}}{{ go_vars $ActionName}}Struct
	// __MYPLUGIN: {{   end }}\
}

// __MYPLUGIN: {{ end }}\

// __MYPLUGIN: {{ range $Objectname, $ObjectOpts := .Yaml.Objects }}\
// __MYPLUGIN: {{   range $ActionName, $Flags := object_tree $ObjectOpts }}\
type AppActionStruct struct {
	// __MYPLUGIN: type {{ go_vars $Objectname}}{{ go_vars $ActionName}}Struct struct {
	// __MYPLUGIN: {{     range $ParamName, $Opts := $Flags }}\
	Param1 string // __MYPLUGIN: 	{{ go_vars $ParamName }} string // {{ $Opts.Help }}
	// __MYPLUGIN: {{     end }}\
}

// __MYPLUGIN: {{   end }}\
// __MYPLUGIN: {{ end }}\

// ************************
// Update request structure
// ************************

type UpdateReq struct {
	Forj struct {
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.common }}\
		Instance string `json:"instance"` // __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		// __MYPLUGIN: {{ end }}\
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.update }}\
		// __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		// __MYPLUGIN: {{ end }}\
	}
	Objects UpdateArgReq
}

type UpdateArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $Opts := .Yaml.Objects }}\
	App map[string]AppInstanceStruct // __MYPLUGIN: 	{{ go_vars $Objectname}} map[string]{{ go_vars $Objectname}}InstanceStruct `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{ end }}\
}

// **************************
// Maintain request structure
// **************************

type MaintainReq struct {
	Forj struct {
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.common }}\
		Instance string `json:"instance"` // __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		// __MYPLUGIN: {{ end }}\
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.maintain }}\
		// __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		// __MYPLUGIN: {{ end }}\
	}
	Objects MaintainArgReq
}

type MaintainArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $Opts := .Yaml.Objects }}\
	// __MYPLUGIN: {{   if object_has_secure $Opts }}\
	App map[string]AppMaintainStruct `json:"app"` // __MYPLUGIN: 	{{ go_vars $Objectname}} map[string]{{ go_vars $Objectname}}MaintainStruct `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{   end }}\
	// __MYPLUGIN: {{ end }}\
}

// __MYPLUGIN: {{ range $Objectname, $ObjectOpts := .Yaml.Objects }}\
// __MYPLUGIN: {{   if object_has_secure $ObjectOpts }}\
type AppMaintainStruct struct {
	// __MYPLUGIN: type {{ go_vars $Objectname}}MaintainStruct struct {
	Setup struct {
		// __MYPLUGIN: {{ range $ParamName, $Opts := $ObjectOpts.Flags }}\
		// __MYPLUGIN: {{   if $Opts.Options.Secure }}\
		Token string // __MYPLUGIN: 		{{ go_vars $ParamName }} string // {{ $Opts.Help }}
		// __MYPLUGIN: {{   end }}\
		// __MYPLUGIN: {{ end }}\
	}
}

// __MYPLUGIN: {{   end }}\
// __MYPLUGIN: {{ end }}\

// YamlDesc has been created from your '{{.Yaml.Name}}.yaml' file.
const YamlDesc = "{{ escape .Yaml_data}}"
