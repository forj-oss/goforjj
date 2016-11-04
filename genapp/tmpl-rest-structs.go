package main

const template_rest_structs = `package main

import "github.hpe.com/christophe-larsonneur/goforjj"

{{ $GroupsList := groups_list .Yaml.Actions }}\
{{ $GroupsMaintainList := groups_list_for "maintain" .Yaml.Actions }}\
{{ if $GroupsList }}\
// Common group of data between create/update actions
{{   range $GroupName, $GroupOpts := $GroupsList }}\
type {{ go_vars $GroupName }}Struct struct {
{{     range $Flagname, $Opts := $GroupOpts.Flags }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{     end }}\
}

{{   end }}\
{{ end }}\
type CreateReq struct {
    Args CreateArgReq ` + "`" + `json:"args"` + "`" + `
    ReposData map[string]goforjj.PluginRepoData
}

type CreateArgReq struct {
{{ if $GroupsList }}\
{{   range $GroupName, $GroupOpts := $GroupsList }}\
    {{ go_vars $GroupName }}Struct
{{   end }}
{{ end }}\
{{ range $Flagname, $Opts := .Yaml.Actions.create.Flags }}\
{{   if $Opts.Group | not }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{   end }}\
{{ end }}\
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
{{   if $Opts.Group | not }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{   end }}\
{{ end }}\
}

type UpdateReq struct {
    Args UpdateArgReq ` + "`" + `json:"args"` + "`" + `
    ReposData map[string]goforjj.PluginRepoData
}

type UpdateArgReq struct {
{{ if $GroupsList }}\
{{   range $GroupName, $GroupOpts := $GroupsList }}\
    {{ go_vars $GroupName }}Struct
{{   end }}
{{ end }}\
{{ range $Flagname, $Opts := .Yaml.Actions.update.Flags }}\
{{   if $Opts.Group | not }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{   end }}\
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
{{   if $Opts.Group | not }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{   end }}\
{{ end }}\
}

type MaintainReq struct {
    Args MaintainArgReq ` + "`" + `json:"args"` + "`" + `
    ReposData map[string]goforjj.PluginRepoData
}

type MaintainArgReq struct {
{{ if $GroupsMaintainList }}\
{{   range $GroupName, $GroupOpts := $GroupsMaintainList }}\
    {{ go_vars $GroupName }}Struct
{{   end }}
{{ end }}\
{{ range $Flagname, $Opts := .Yaml.Actions.maintain.Flags }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
{{   if $Opts.Group | not }}\
    {{ go_vars $Flagname}} string ` + "`" + `json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{   end }}\
{{ end }}\
}

// Function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *CreateArgReq)SaveMaintainOptions(ret *goforjj.PluginData) {
    if ret.Options == nil {
        ret.Options = make(map[string]goforjj.PluginOption)
    }
{{ range $Flagname, $Opts := maintain_options "create" .Yaml.Actions }}
{{   if has_prefix $Flagname "forjj-" | not }}\
    ret.Options["{{ $Flagname }}"] = addMaintainOptionValue(ret.Options, "{{ $Flagname }}", r.{{ go_vars $Flagname }}, "{{ $Opts.Default }}", "{{ $Opts.Help }}")
{{   end }}\
{{ end }}\
}

func (r *UpdateArgReq)SaveMaintainOptions(ret *goforjj.PluginData) {
    if ret.Options == nil {
        ret.Options = make(map[string]goforjj.PluginOption)
    }
{{ range $Flagname, $Opts := maintain_options "update" .Yaml.Actions }}
{{   if has_prefix $Flagname "forjj-" | not }}\
    ret.Options["{{ $Flagname }}"] = addMaintainOptionValue(ret.Options, "{{ $Flagname }}", r.{{ go_vars $Flagname }}, "{{ $Opts.Default }}", "{{ $Opts.Help }}")
{{   end }}\
{{ end }}\
}

func addMaintainOptionValue(options map[string]goforjj.PluginOption, option, value, defaultv, help string) (goforjj.PluginOption){
    opt, ok := options[option]
    if ok && value != "" {
        opt.Value = value
        return opt
    }
    if ! ok {
        opt = goforjj.PluginOption { Help: help }
        if value == "" {
            opt.Value = defaultv
        } else {
            opt.Value = value
        }
    }
    return opt
}

// YamlDesc has been created from your '{{.Yaml.Name}}.yaml' file.
const YamlDesc="{{ escape .Yaml_data}}"

`
