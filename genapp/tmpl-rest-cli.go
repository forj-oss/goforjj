package main

const template_rest_cli = `package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/forjj-oss/goforjj"
    "gopkg.in/yaml.v2"
    "log"
)

type {{.Yaml.Name}}App struct {
    App *kingpin.Application
    params Params
    socket string
    Yaml goforjj.YamlPlugin
}

type Params struct {
    socket_file *string
    socket_path *string
    daemon *bool // Currently not used - Lot of concerns with daemonize in go... Stay in foreground
}

func (a *{{.Yaml.Name}}App)init() {
    a.load_plugin_def()

    a.App = kingpin.New("{{.Yaml.Name}}", "{{.Yaml.Description}}")
{{ if .Yaml.Version }}\
    a.App.Version("{{ .Yaml.Version }}")
{{ end }}\

    // true to create the Infra
    daemon := a.App.Command("service", "{{ .Yaml.Name}} REST API service")
    daemon.Command("start", "start {{ .Yaml.Name}} REST API service")
    a.params.socket_file = daemon.Flag("socket-file", "Socket file to use").Default(a.Yaml.Runtime.Service.Socket).String()
    a.params.socket_path = daemon.Flag("socket-path", "Socket file path to use").Default("/tmp/forjj-socks").String()
    a.params.daemon = daemon.Flag("daemon", "Start process in background like a daemon").Short('d').Bool()
}

func (a *{{.Yaml.Name}}App)load_plugin_def() {
    yaml.Unmarshal([]byte(YamlDesc), &a.Yaml)
    if a.Yaml.Runtime.Service.Socket == "" {
        a.Yaml.Runtime.Service.Socket = "{{ .Yaml.Name }}.sock"
        log.Printf("Set default socket file: %s", a.Yaml.Runtime.Service.Socket)
    }
}
`

/*
  a.Flags = make(map[string]*string)
  a.Tasks = make(map[string]goforjj.PluginTask)
{{ range $Flagname, $Opts := .Actions.common.Flags }}\
  a.Flags["{{ $Flagname }}"] = a.App.Flag("{{ $Flagname }}", "{{ $Opts.Help }}")\
{{   if $Opts.Required }}.Required(){{ end }}\
{{   if $Opts.Default  }}.Default("{{$Opts.Default}}"){{ end }}\
{{   if $Opts.Hidden   }}.Hidden(){{ end }}.String()
{{ end }}

{{ range $cmd, $flags := (filter_cmds .Actions) }}\
  a.Tasks["{{$cmd}}"] = goforjj.PluginTask {
    Flags: make(map[string]*string),
    Cmd  : a.App.Command("{{ $cmd }}", "{{ $flags.Help }}"),
  }
{{   range $Flagname, $Opts := $flags.Flags }}\
  a.Tasks["{{$cmd}}"].Flags["{{ $Flagname }}"] = a.Tasks["{{$cmd}}"].Cmd.Flag("{{ $Flagname }}", "{{ $Opts.Help }}")\
{{     if $Opts.Required }}.Required(){{ end }}\
{{     if $Opts.Default  }}.Default("{{$Opts.Default}}"){{ end }}\
{{     if $Opts.Hidden   }}.Hidden(){{ end }}.String()
{{   end }}
{{ end }}\
*/
