package main

const template_rest_maintain = `package main

import (
    "github.com/forjj-oss/goforjj"
    "log"
    "os"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) check_source_existence(ret *goforjj.PluginData) (status bool) {
    log.Printf("Checking {{ go_vars .Yaml.Name }} source code path existence.")

    if _, err := os.Stat(r.ForjjSourceMount) ; err == nil {
        ret.Errorf("Unable to maintain {{ .Yaml.Name }} instances. '%s' is inexistent or innacessible.\n", r.ForjjSourceMount)
        return
    }
    ret.StatusAdd("environment checked.")
    return true
}

func (r *MaintainReq)instantiate(ret *goforjj.PluginData) (status bool) {

    return true
}
`
