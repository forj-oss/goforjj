package main

const template_rest_structs = `package main

type CreateReq struct {
{{ range $Flagname, $Opts := .Yaml.Actions.create.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`yaml:"{{$Flagname}}"` + "`" + ` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`yaml:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}\
}

type UpdateReq struct {
{{ range $Flagname, $Opts := .Yaml.Actions.update.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`yaml:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`yaml:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}\
}

type MaintainReq struct {
{{ range $Flagname, $Opts := .Yaml.Actions.maintain.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`yaml:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}
    // common flags
{{ range $Flagname, $Opts := .Yaml.Actions.common.Flags }}\
    {{ go_vars $Flagname}} string `+"`"+`yaml:"{{$Flagname}}"`+"`"+` // {{ $Opts.Help }}
{{ end }}\
}

// YamlDesc has been created from your '{{.Yaml.Name}}.yaml' file.
const YamlDesc="{{ escape .Yaml_data}}"

`
