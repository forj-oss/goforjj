package main

import (
	"github.com/forj-oss/goforjj"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

type __MYPLUGIN__Plugin struct {
	yaml        Yaml__MYPLUGIN__
	source_path string
}

const __MYPLUGIN_UNDERSCORED___file = "forjj-__MYPLUGINNAME__.yaml"

type Yaml__MYPLUGIN__ struct {
}

func new_plugin(src string) (p *__MYPLUGIN__Plugin) {
	p = new(__MYPLUGIN__Plugin)

	p.source_path = src
	return
}

func (p *__MYPLUGIN__Plugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *__MYPLUGIN__Plugin) load_from(ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *__MYPLUGIN__Plugin) update_from(r *UpdateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *__MYPLUGIN__Plugin) save_yaml(ret *goforjj.PluginData, instance string) (status bool) {
	file := path.Join(instance, __MYPLUGIN_UNDERSCORED___file)

	d, err := yaml.Marshal(&p.yaml)
	if err != nil {
		ret.Errorf("Unable to encode forjj __MYPLUGINNAME__ configuration data in yaml. %s", err)
		return
	}

	if err := ioutil.WriteFile(file, d, 0644); err != nil {
		ret.Errorf("Unable to save '%s'. %s", file, err)
		return
	}
	return true
}

func (p *__MYPLUGIN__Plugin) load_yaml(ret *goforjj.PluginData, instance string) (status bool) {
	file := path.Join(instance, __MYPLUGIN_UNDERSCORED___file)

	d, err := ioutil.ReadFile(file)
	if err != nil {
		ret.Errorf("Unable to load '%s'. %s", file, err)
		return
	}

	err = yaml.Unmarshal(d, &p.yaml)
	if err != nil {
		ret.Errorf("Unable to decode forjj __MYPLUGINNAME__ data in yaml. %s", err)
		return
	}
	return true
}
