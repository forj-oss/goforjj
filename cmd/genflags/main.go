package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "os"
    //"os/exec"
    "io/ioutil"
    //"io"
    "strings"
    "text/template"
)

type Flags struct {
    Opts_m map[string]string
}

type Cmds struct {
    Help    string
    Flags_m map[string]Flags
}

type Appflags struct {
    YamlDesc string
    Plugin   string
    Version  string
    Desc     string
    Cmds_m   map[string]Cmds
}

func (a *Appflags) ReadFrom(data interface{}) {
    a.Cmds_m = make(map[string]Cmds)

    for root_i, i := range data.(map[interface{}]interface{}) {
        root := root_i.(string)

        switch root {
        case "flags":
            a.ReadFlagsFrom(i)
        case "plugin":
            a.Plugin = i.(string)
        case "version":
            a.Version = i.(string)
        case "description":
            a.Desc = i.(string)
        }
    }
}

func (a *Appflags) ReadFlagsFrom(i interface{}) {

    cmd_flags := i.(map[interface{}]interface{})
    for cmd_i, i := range cmd_flags {
        cmd := cmd_i.(string)
        a.Cmds_m[cmd] = Cmds{Flags_m: make(map[string]Flags)}
        cmd_s := a.Cmds_m[cmd]

        if i == nil {
            continue
        }

        Cmds_opts := i.(map[interface{}]interface{})
        for cmd_opt_i, i := range Cmds_opts {
            cmd_opt := cmd_opt_i.(string)
            switch cmd_opt {
            case "help":
                cmd_s.Help = i.(string)
            case "flags":
                cmd_s.ReadCmdFlagsFrom(i)
            }
        }
    }
}

func (c *Cmds) ReadCmdFlagsFrom(i interface{}) {
    flags_opts := i.([]interface{})
    for _, i := range flags_opts {
        if i == nil {
            continue
        }

        flag_opts := i.(map[interface{}]interface{})
        for flag_i, i := range flag_opts {

            flag := flag_i.(string)
            c.Flags_m[flag] = Flags{make(map[string]string)}

            if i == nil {
                continue
            }

            opts_i := i.(map[interface{}]interface{})
            for k, v := range opts_i {
                switch v.(type) {
                case string:
                    c.Flags_m[flag].Opts_m[k.(string)] = v.(string)
                case bool:
                    if v.(bool) {
                        c.Flags_m[flag].Opts_m[k.(string)] = "true"
                    }
                }
            }
        }
    }
}

func main() {
    var (
        yaml_data   []byte
        data        interface{}
        tmpl        *template.Template
        flags       Appflags
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

    if err := yaml.Unmarshal(yaml_data, &data); err != nil {
        fmt.Printf("go-forjj-generate: error! '%s' is not a valid yaml document. %s\n", source_file, err)
        os.Exit(1)
    }

    flags.ReadFrom(data)
    flags.YamlDesc = string(yaml_data)

    if flags.Plugin == "" {
        fmt.Printf("go-forjj-generate: error! '%s' missed '/plugin' key.\n", source_file)
        os.Exit(1)
    }

    if t, err := template.New("genflags").Funcs(template.FuncMap{
        "escape": func(str string) string {
            return strings.Replace(strings.Replace(str, "\"", "\\\"", -1), "\n", "\\n\" +\n   \"", -1)
        },
        "Test": func(opts map[string]string, option, value string) string {
            if v, ok := opts[option]; ok && v == value {
                return "true"
            }
            return ""
        },
        "filter_cmds": func(flags map[string]Cmds) (ret map[string]Cmds) {
            ret = make(map[string]Cmds)
            for _, v := range []string{"check", "create", "update", "maintain"} {
                ret[v] = Cmds{}
            }
            for k, v := range flags {
                if k != "common" {
                    ret[k] = v
                }
            }
            return
        },
    }).Parse(strings.Replace(template_source, "\\\n", "", -1)); err != nil {
        fmt.Printf("go-forjj-generate: Template error: %s\n", err)
        os.Exit(1)
    } else {
        tmpl = t
    }

    w, err := os.Create(flags.Plugin + ".go")
    if err != nil {
        fmt.Printf("go-forjj-generate: error! '%s' is not writeable. %s\n", source_file, err)
        os.Exit(1)
    }
    defer w.Close()

    if err := tmpl.Execute(w, flags); err != nil {
        fmt.Printf("go-forjj-generate: error! %s\n", err)
        os.Exit(1)
    }
    //if err := exec.Command("goimports", "-w", "flags_generated.go").Run() ; err != nil {
    //}
}
