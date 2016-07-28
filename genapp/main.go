package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "os"
    "io/ioutil"
    "github.hpe.com/christophe-larsonneur/goforjj"
)

type App struct {
    Yaml goforjj.YamlPlugin
    Models
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

    if os.Args[1] == "" {
        fmt.Printf("go-forjj-generate: Error! Yaml source file missing.\n")
        os.Exit(1)
    }

    source_file = os.Args[1]
    if _, err := os.Stat(source_file); os.IsNotExist(err) {
        fmt.Printf("go-forjj-generate: Error! Yaml source file '%s' is not accessible.\n", source_file)
        os.Exit(1)
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

    app.init_model()

    app.Models.Create_model(&app.Yaml, yaml_data, app.Yaml.Runtime.Service_type)

}
