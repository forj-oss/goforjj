package main

const template_rest_update = `package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "path"
    "os"
)

// Return ok if the jenkins instance exist
func (r *UpdateReq) check_source_existence(ret *goforjj.PluginData) (p *{{ go_vars .Yaml.Name }}Plugin, status bool) {
    log.Printf("Checking {{ go_vars .Yaml.Name }} source code existence.")
    src_path := path.Join(r.ForjjSourceMount, r.ForjjInstanceName)
    if _, err := os.Stat(path.Join(src_path, {{ go_vars_underscored .Yaml.Name }}_file)) ; err == nil {
        log.Printf(ret.Errorf("Unable to create the {{ .Yaml.Name }} source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update {{ .Yaml.Name }} according to his configuration. %s.", src_path, src_path, err))
        return
    }

    p = new_plugin(src_path)

    ret.StatusAdd("environment checked.")
    return p, true
}

func (r *{{ go_vars .Yaml.Name }}Plugin)update_jenkins_sources(ret *goforjj.PluginData) (status bool) {
    return true
}
`
