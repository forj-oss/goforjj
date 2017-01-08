package main

import (
	"fmt"
	"github.com/forjj-oss/goforjj"
	"os"
	"path"
	"strings"
	"text/template"
)

const prefix_generated_template = `// This file is autogenerated by "go generate". Do not modify it.
// It has been generated from your '{{.Yaml.Name}}.yaml' file.
// To update those structure, update the '{{.Yaml.Name}}.yaml' and run 'go generate'
`

const prefix_created_template = `// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
`

type YamlData struct {
	Yaml      *goforjj.YamlPlugin
	Yaml_data string
}

// Define a source code model. Depending on the plugin service type, the plugin initial sources will be created from a model of sources (REST API or shell for example)
type Models struct {
	model map[string]Model
}

// Collection of source files to generate
type Model struct {
	sources map[string]Source // key is the filename
}

// Core of the generated source file
// If reset = true, the generated file will regenerated each go build done.
type Source struct {
	reset    bool
	template string
	rights   os.FileMode
}

// Create a new model of sources
func (m *Models) Create(name string) *Model {
	if m.model == nil {
		m.model = make(map[string]Model)
	}

	sources := Model{make(map[string]Source)}
	m.model[name] = sources
	return &sources
}

// Add a source template to the model
func (m *Model) Source(file string, rights os.FileMode, comment, tmpl_src string, reset bool) *Model {
	template := strings.Replace(tmpl_src, "\\\n", "", -1)
	source := Source{}
	source.reset = reset
	source.rights = rights

	switch {
	case comment == "":
		source.template = template
	case reset:
		file = "generated-" + file
		source.template = template_comment(prefix_generated_template, comment) + template
	case !reset:
		source.template = template_comment(prefix_created_template, comment) + template
	}

	m.sources[file] = source
	return m
}

// Set appropriate comment prefix
func template_comment(template, comment string) string {
	return strings.Replace(template, "//", comment, -1)
}

// Create the source files from the model given.
func (m *Models) Create_model(yaml *goforjj.YamlPlugin, raw_yaml []byte, name string) {
	var yaml_data YamlData = YamlData{yaml, string(raw_yaml)}

	model, ok := m.model[name]
	if !ok {
		fmt.Printf("Invalid Model '%s' to apply.\n", name)
		os.Exit(1)
	}
	for k, v := range model.sources {
		v.apply_source(&yaml_data, k)
	}
}

func (s *Source) apply_source(yaml *YamlData, file string) {
	var tmpl *template.Template
	var err error
	var w *os.File

	if _, err = os.Stat(file); err == nil && !s.reset {
		return
	}

	// TODO: Normalize Structure name in template. For ex, - is not supported. Replace it to _ or remove it.
	tmpl, err = template.New(file).Funcs(template.FuncMap{
		"escape": func(str string) string {
			return strings.Replace(strings.Replace(str, "\"", "\\\"", -1), "\n", "\\n\" +\n   \"", -1)
		},
		"go_vars": func(str string) string {
			return strings.Replace(strings.Title(str), "-", "", -1)
		},
		"go_vars_underscored": func(str string) string {
			return strings.Replace(str, "-", "_", -1)
		},
		"groups_list": func(actions map[string]goforjj.YamlPluginDef) (ret map[string]goforjj.YamlPluginDef) {
			ret = make(map[string]goforjj.YamlPluginDef)
			for name, action_opts := range actions {
				if name != "create" && name != "update" && name != "common" {
					continue
				}
				for flag_name, flag_opts := range action_opts.Flags {
					if flag_opts.Group == "" {
						continue
					}
					if group, found := ret[flag_opts.Group]; found {
						if _, found := group.Flags[flag_name]; found {
							continue
						}
						group.Flags[flag_name] = flag_opts
						ret[flag_opts.Group] = group
					} else {
						group := goforjj.YamlPluginDef{Flags: make(map[string]goforjj.YamlFlagsOptions)}
						group.Flags[flag_name] = flag_opts
						ret[flag_opts.Group] = group
					}
				}
			}
			return
		},
		"groups_list_for": func(cmd string, actions map[string]goforjj.YamlPluginDef) (ret map[string]goforjj.YamlPluginDef) {
			ret = make(map[string]goforjj.YamlPluginDef)
			for name, action_opts := range actions {
				if name != cmd && name != "common" {
					continue
				}
				for flag_name, flag_opts := range action_opts.Flags {
					if flag_opts.Group == "" {
						continue
					}
					if group, found := ret[flag_opts.Group]; found {
						if _, found := group.Flags[flag_name]; found {
							continue
						}
						group.Flags[flag_name] = flag_opts
						ret[flag_opts.Group] = group
					} else {
						group := goforjj.YamlPluginDef{Flags: make(map[string]goforjj.YamlFlagsOptions)}
						group.Flags[flag_name] = flag_opts
						ret[flag_opts.Group] = group
					}
				}
			}
			return
		},
		"maintain_options": func(action string, actions map[string]goforjj.YamlPluginDef) (ret map[string]goforjj.YamlFlagsOptions) {
			ret = make(map[string]goforjj.YamlFlagsOptions)
			// Get a list of secure values defined in create/update phase
			for ak, av := range actions {
				if ak != "maintain" {
					continue
				}
				// Check each maintain flags exist on 'create/update' list of flags.
				for fn, fv := range av.Flags {
					if fd, ok := actions[action].Flags[fn]; ok && fd.Secure {
						ret[fn] = fv
					}
					if fd, ok := actions["common"].Flags[fn]; ok && fd.Secure {
						ret[fn] = fv
					}
				}
			}
			// All common case identified secure are added as well.
			for fn, fv := range actions["common"].Flags {
				if fv.Secure {
					ret[fn] = fv
				}

			}
			return
		},
		"has_prefix": strings.HasPrefix,
		"filter_cmds": func(actions map[string]goforjj.YamlPluginDef) (ret map[string]goforjj.YamlPluginDef) {
			ret = make(map[string]goforjj.YamlPluginDef)

			for k, v := range actions {
				if k != "common" {
					ret[k] = v
				}
			}
			return
		},
	}).Parse(s.template)
	if err != nil {
		fmt.Printf("go-forjj-generate: Template error: %s\n", err)
		os.Exit(1)
	}

	file_path := path.Dir(file)

	if file_path != "." {
		if err := os.MkdirAll(file_path, 0755); err != nil {
			fmt.Printf("go-forjj-generate: error! Unable to create '%s' tree\n", file_path)
			os.Exit(1)
		}
	}

	w, err = os.Create(file)
	if err != nil {
		fmt.Printf("go-forjj-generate: error! '%s' is not writeable. %s\n", file, err)
		os.Exit(1)
	}
	defer w.Close()

	if err = tmpl.Execute(w, yaml); err != nil {
		fmt.Printf("go-forjj-generate: error! %s\n", err)
		os.Exit(1)
	}

	if err := os.Chmod(file, s.rights); err != nil {
		fmt.Printf("go-forjj-generate: error! Unable to set rights %d. %s\n", s.rights, err)
		os.Exit(1)
	}

	if s.reset {
		fmt.Printf("%s\n", file)
	} else {
		fmt.Printf("'%s' created. Won't be updated anymore at next go generate until file disappear.\n", file)
	}
}

const yaml_template = `---
plugin: "{{ .Yaml.Name }}"
version: "0.1"
description: "{{ .Yaml.Name }} plugin for FORJJ."
runtime:
  docker_image: "docker.hos.hpecorp.net/forjj/{{ .Yaml.Name }}"
  service_type: "REST API"
  service:
    #socket: "{{ .Yaml.Name }}.sock"
    parameters: [ "service", "start" ]
created_flag_file: "{{ "{{ .InstanceName }}" }}/forjj-{{ "{{ .Name }}" }}.yaml"
actions:
  common:
    flags:
      forjj-infra:
        help: "Name of the Infra repository to use"
      {{ .Yaml.Name }}-debug:
        help: "To activate {{ .Yaml.Name }} debug information"
      forjj-source-mount: # Used by the plugin to store plugin data in yaml. See {{ go_vars_underscored .Yaml.Name }}_plugin.go
        help: "Where the source dir is located for {{ .Yaml.Name }} plugin."
  create:
    help: "Create a {{ .Yaml.Name }} instance source code."
    flags:
      # Options related to source code
      forjj-instance-name: # Used by the plugin to store plugin data in yaml for the current instance. See {{ go_vars_underscored .Yaml.Name }}_plugin.go
        help: "Name of the {{ .Yaml.Name }} instance given by forjj."
        group: "source"
  update:
    help: "Update a {{ .Yaml.Name }} instance source code"
    flags:
      forjj-instance-name: # Used by the plugin to store plugin data in yaml for the current instance. See {{ go_vars_underscored .Yaml.Name }}_plugin.go
        help: "Name of the {{ .Yaml.Name }} instance given by forjj."
        group: "source"
  maintain:
    help: "Instantiate {{ .Yaml.Name }} thanks to source code."
`
