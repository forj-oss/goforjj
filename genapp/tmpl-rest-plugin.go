package main

const template_rest_plugin = `package main

import (
    "github.com/forj-oss/goforjj"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "path"
)

type {{ go_vars .Yaml.Name }}Plugin struct {
    yaml Yaml{{ go_vars .Yaml.Name }}
    source_path string
}

const {{ go_vars_underscored .Yaml.Name }}_file = "forjj-{{ .Yaml.Name }}.yaml"

type Yaml{{ go_vars .Yaml.Name }} struct {
{{ $GroupsList := groups_list .Yaml.Actions }}\
{{ if $GroupsList }}\
{{   range $GroupName, $GroupOpts := $GroupsList }}\
    {{ go_vars $GroupName }} {{ go_vars $GroupName }}Struct
{{   end}}\
{{ else}}\
    // Given as an example:
    data string
{{ end}}\
}

func new_plugin(src string) (p *{{ go_vars .Yaml.Name }}Plugin) {
    p = new({{ go_vars .Yaml.Name }}Plugin)

    p.source_path = src
    return
}

func (p *{{ go_vars .Yaml.Name }}Plugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
{{ $GroupsList := groups_list .Yaml.Actions }}\
{{ if $GroupsList }}\
{{   range $GroupName, $GroupOpts := $GroupsList }}\
    p.yaml.{{ go_vars $GroupName }} = r.{{ go_vars $GroupName }}Struct
{{   end}}\
{{ else}}\
    // Given as an example:
    data = "example"
{{ end}}\
    return true
}

func (p *{{ go_vars .Yaml.Name }}Plugin) load_from(ret *goforjj.PluginData) (status bool) {
    return true
}

func (p *{{ go_vars .Yaml.Name }}Plugin) update_from(r *UpdateReq, ret *goforjj.PluginData)  (status bool) {
    return true
}

func (p *{{ go_vars .Yaml.Name }}Plugin)save_yaml(ret *goforjj.PluginData) (status bool) {
    file := path.Join(p.yaml.Source.ForjjInstanceName, {{ go_vars_underscored .Yaml.Name }}_file)

    d, err := yaml.Marshal(&p.yaml)
    if  err != nil {
        ret.Errorf("Unable to encode forjj {{ .Yaml.Name }} configuration data in yaml. %s", err)
        return
    }

    if err := ioutil.WriteFile(file, d, 0644) ; err != nil {
        ret.Errorf("Unable to save '%s'. %s", file, err)
        return
    }
    return true
}

func (p *{{ go_vars .Yaml.Name }}Plugin)load_yaml(ret *goforjj.PluginData) (status bool) {
    file := path.Join(p.yaml.Source.ForjjInstanceName, {{ go_vars_underscored .Yaml.Name }}_file)

    d, err := ioutil.ReadFile(file)
    if err != nil {
        ret.Errorf("Unable to load '%s'. %s", file, err)
        return
    }

    err = yaml.Unmarshal(d, &p.yaml)
    if  err != nil {
        ret.Errorf("Unable to decode forjj {{ .Yaml.Name }} data in yaml. %s", err)
        return
    }
    return true
}
`
