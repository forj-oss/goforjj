package main

// __MYPLUGIN: {{ range $ObjectName, $Object := .Yaml.Objects }}\
// Object {{ $ObjectName }} groups structure

// __MYPLUGIN: {{   range $GroupName, $Group := $Object.Groups }}\
type DataStruct struct { // __MYPLUGIN: type {{ go_vars $GroupName }}Struct struct {
	// __MYPLUGIN: {{     range $FlagName, $Flag := $Group.Flags }}\
	data1 string // __MYPLUGIN: 	{{ go_vars $FlagName }} string
	// __MYPLUGIN: {{     end }}\
}

// Action groups structure
// __MYPLUGIN: {{   end }}\
// __MYPLUGIN: {{   range $ActionName, $Action := object_tree $Object }}\
// __MYPLUGIN: {{     range $GroupName, $Group := $Action.Groups }}
type ActionDataStruct struct { // __MYPLUGIN: type {{go_vars $ActionName}}{{ go_vars $GroupName }}Struct struct {
	// __MYPLUGIN: {{       range $FlagName, $Flag := $Group.Flags }}\
	// __MYPLUGIN: {{         if $Flag.Actions }}
	// __MYPLUGIN: {{           if inList $ActionName $Flag.Actions }}\
	data1 string // __MYPLUGIN: 	{{ go_vars $FlagName }} string
	// __MYPLUGIN: {{           end }}\
	// __MYPLUGIN: {{         else }}\
	// __MYPLUGIN: 	{{ go_vars $FlagName }} string
	// __MYPLUGIN: {{         end }}\
	// __MYPLUGIN: {{       end }}\
}

// __MYPLUGIN: {{     end }}\
// __MYPLUGIN: {{   end }}\

// Object Instance structures

type AppInstanceStruct struct { // __MYPLUGIN: type {{ go_vars $ObjectName}}InstanceStruct struct {
	// __MYPLUGIN: {{   range $ActionName, $Flags := object_tree $Object }}\
	Action AppActionStruct // __MYPLUGIN: 	{{ go_vars $ActionName}} {{ go_vars $ObjectName}}{{ go_vars $ActionName}}Struct
	// __MYPLUGIN: {{   end }}\
}

// Object instance Action structures

// __MYPLUGIN: {{   range $ActionName, $Action := object_tree $Object }}\
type AppActionStruct struct { // __MYPLUGIN: type {{ go_vars $ObjectName}}{{ go_vars $ActionName}}Struct struct {
	// __MYPLUGIN: {{     range $ParamName, $Opts := $Action.Flags }}\
	Param1 string // __MYPLUGIN: 	{{ go_vars $ParamName }} string // {{ $Opts.Help }}
	// __MYPLUGIN: {{     end }}\
	// __MYPLUGIN: {{     if $Action.Groups }}\

	// Groups

	// __MYPLUGIN: {{     end }}\
	// __MYPLUGIN: {{     range $GroupName, $Group := $Action.Groups }}\
	ActionDataStruct // __MYPLUGIN: 	{{go_vars $ActionName}}{{ go_vars $GroupName }}Struct
	// __MYPLUGIN: {{     end }}\
}

// __MYPLUGIN: {{   end }}\
// __MYPLUGIN: {{ end }}\

// ************************
// Create request structure
// ************************

type CreateReq struct {
	Forj struct {
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.common }}\
		ForjjInstanceName string `json:"forjj-instance-name"` // __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		ForjjSourceMount  string `json:"forjj-source-mount"`  // __MYPLUGIN: {{ end }}\
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

// ************************
// Update request structure
// ************************

type UpdateReq struct {
	Forj struct {
		// __MYPLUGIN: {{ range $FlagName, $FlagOpts := .Yaml.Tasks.common }}\
		ForjjInstanceName string `json:"forjj-instance-name"` // __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		ForjjSourceMount  string `json:"forjj-source-mount"`  // __MYPLUGIN: {{ end }}\
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
		ForjjInstanceName string `json:"forjj-instance-name"` // __MYPLUGIN: 		{{ go_vars $FlagName }} string `json:"{{ $FlagName }}"`
		ForjjSourceMount  string `json:"forjj-source-mount"`  // __MYPLUGIN: {{ end }}\
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
type AppMaintainStruct struct { // __MYPLUGIN: type {{ go_vars $Objectname}}MaintainStruct struct {
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
