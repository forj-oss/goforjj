package main

const template_rest_structs = `package main

import "github.hpe.com/christophe-larsonneur/goforjj"

type CreateReq struct {
{{ range $Flagname, $Opts := .Yaml.Actions.create.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`json:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`json:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}\
}

type UpdateReq struct {
{{ range $Flagname, $Opts := .Yaml.Actions.update.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`json:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`json:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}\
}

type MaintainReq struct {
{{ range $Flagname, $Opts := .Yaml.Actions.maintain.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`json:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`json:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}\
}

// Function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *CreateReq)SaveMaintainOptions(ret *goforjj.PluginData) {
    if ret.Options == nil {
        ret.Options = make(map[string]goforjj.PluginOption)
    }
{{ range $Flagname, $Opts := maintain_options "create" .Yaml.Actions }}
{{   if has_prefix $Flagname "forjj-" | not }}\
    ret.Options["{{ $Flagname }}"] = addMaintainOptionValue(ret.Options, "{{ $Flagname }}", r.{{ go_vars $Flagname }}, "{{ $Opts.Default }}", "{{ $Opts.Help }}")
{{   end }}\
{{ end }}\
}

func (r *UpdateReq)SaveMaintainOptions(ret *goforjj.PluginData) {
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
