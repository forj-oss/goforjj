package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *__MYPLUGIN__Plugin, status bool) {
	log.Print("Checking __MYPLUGIN__ source code existence.")
	src_path := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(src_path, __MYPLUGIN_UNDERSCORED___file)); err == nil {
		log.Printf(ret.Errorf("Unable to create the __MYPLUGINNAME__ source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update __MYPLUGINNAME__ according to his configuration. %s.", src_path, src_path, err))
		return
	}

	p = new_plugin(src_path)

	log.Printf(ret.StatusAdd("environment checked."))
	return p, true
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}
