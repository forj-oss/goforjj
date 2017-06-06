package main

import (
	"fmt"
	"goforjj"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type App struct {
	Yaml goforjj.YamlPlugin
	Models
	template_path string
}

func (a *App) ReadFrom(yaml_data []byte) error {
	return yaml.Unmarshal(yaml_data, &a.Yaml)
}

func main() {
	var (
		yaml_data   []byte
		app         App
		source_file string
	)

	if len(os.Args) != 3 {
		fmt.Print("Usage is : genapp <plugin yaml file> <template source path>")
		os.Exit(1)
	}
	if os.Args[1] == "" {
		fmt.Print("go-forjj-generate: Error! Yaml source file missing.\n")
		os.Exit(1)
	}
	if os.Args[2] == "" {
		fmt.Print("go-forjj-generate: Error! template source files missing." +
			"It should be something like '.../goforjj/genapp'\n")
		os.Exit(1)
	}

	source_file = os.Args[1]
	app.template_path = os.Args[2]
	if _, err := os.Stat(app.template_path); os.IsNotExist(err) {
		fmt.Printf("go-forjj-generate: Warning! template source files path '%s' is not accessible. "+
			"It should be something like '.../goforjj/genapp'\n", app.template_path)
		os.Exit(1)
	}

	if _, err := os.Stat(source_file); os.IsNotExist(err) {
		fmt.Printf("go-forjj-generate: Warning! Yaml source file '%s' is not accessible. Trying to create a basic one\n", source_file)

		app.Yaml.Name = strings.Replace(source_file, ".yaml", "", -1)
		yaml_source := Source{
			template: yaml_template,
			reset:    false,
			rights:   0644,
		}

		yaml_data := &YamlData{
			Yaml: &app.Yaml,
		}

		yaml_source.apply_source(yaml_data, source_file)
		fmt.Printf("go-forjj-generate: Reading the Yaml source file '%s' created\n", source_file)
	}

	if d, err := ioutil.ReadFile(source_file); err != nil {
		fmt.Printf("go-forjj-generate: Error! '%s' is not a readable document. %s\n", source_file, err)
		os.Exit(1)
	} else {
		yaml_data = d
	}

	if err := app.ReadFrom(yaml_data); err != nil {
		fmt.Printf("go-forjj-generate: error! '%s' is not a valid yaml document. %s\n", source_file, err)
		os.Exit(1)
	}

	if app.Yaml.Name == "" {
		fmt.Printf("go-forjj-generate: error! '%s' missed '/plugin' key.\n", source_file)
		os.Exit(1)
	}
	if app.Yaml.Runtime.Service_type == "" {
		fmt.Printf("go-forjj-generate: error! '%s' missed '/Runtime/Service_type'. Valid values 'REST API', 'shell'\n", source_file)
		os.Exit(1)
	}

	model_id := app.init_model()

	app.Models.Create_model(&app.Yaml, yaml_data, model_id)

}
