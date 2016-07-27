package main

import (
    "strings"
    "text/template"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "fmt"
    "os"
)

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
    reset bool
    template string
}

// Create a new model of sources
func (m *Models)Create(name string) (*Model) {
    if m.model == nil {
        m.model = make(map[string]Model)
    }

    sources := Model {make(map[string]Source)}
    m.model[name] = sources
    return &sources
}

// Add a source template to the model
func (m *Model)Source(file, tmpl_src string, reset bool) (*Model){
    source := Source{}
    source.reset = reset
    source.template = strings.Replace(tmpl_src, "\\\n", "", -1)
    m.sources[file] = source
    return m
}

// Create the source files from the model given.
func (m *Models)Create_model(yaml *goforjj.YamlPlugin, name string) {
    model, ok := m.model[name]
    if !ok {
        fmt.Printf("Invalid Model '%s' to apply.\n", name)
        os.Exit(1)
    }
    for k, v := range model.sources {
        v.apply_source(yaml, k)
    }
}

func (s *Source)apply_source(yaml *goforjj.YamlPlugin, file string) {
    var tmpl *template.Template
    var err error
    var w *os.File

    tmpl, err = template.New(file).Funcs(template.FuncMap{
        "escape": func(str string) string {
            return strings.Replace(strings.Replace(str, "\"", "\\\"", -1), "\n", "\\n\" +\n   \"", -1)
        },
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

}
