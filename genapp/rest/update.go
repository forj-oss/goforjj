package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

// Return ok if the jenkins instance exist
func (r *UpdateReq) check_source_existence(ret *goforjj.PluginData) (p *__MYPLUGIN__Plugin, status bool) {
	log.Print("Checking __MYPLUGIN__ source code existence.")
	src_path := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(src_path, __MYPLUGIN_UNDERSCORED___file)); err == nil {
		log.Printf(ret.Errorf("Unable to create the __MYPLUGINNAME__ source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update __MYPLUGINNAME__ according to his configuration. %s.", src_path, src_path, err))
		return
	}

	p = new_plugin(src_path)

	ret.StatusAdd("environment checked.")
	return p, true
}

func (r *__MYPLUGIN__Plugin) update_jenkins_sources(ret *goforjj.PluginData) (status bool) {
	return true
}

// Function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *UpdateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

func addMaintainOptionValue(options map[string]goforjj.PluginOption, option, value, defaultv, help string) goforjj.PluginOption {
	opt, ok := options[option]
	if ok && value != "" {
		opt.Value = value
		return opt
	}
	if !ok {
		opt = goforjj.PluginOption{Help: help}
		if value == "" {
			opt.Value = defaultv
		} else {
			opt.Value = value
		}
	}
	return opt
}
