package main

const template_rest_create = `package main

import (
    "github.com/forj-oss/goforjj"
    "os"
    "log"
    "path"
)

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *{{ go_vars .Yaml.Name }}Plugin, status bool) {
    log.Printf("Checking {{ go_vars .Yaml.Name }} source code existence.")
    src_path := path.Join(r.ForjjSourceMount, r.ForjjInstanceName)
    if _, err := os.Stat(path.Join(src_path, {{ go_vars_underscored .Yaml.Name }}_file)) ; err == nil {
        log.Printf(ret.Errorf("Unable to create the {{ .Yaml.Name }} source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update {{ .Yaml.Name }} according to his configuration. %s.", src_path, src_path, err))
        return
    }

    p = new_plugin(src_path)

    log.Printf(ret.StatusAdd("environment checked."))
    return p, true
}`
