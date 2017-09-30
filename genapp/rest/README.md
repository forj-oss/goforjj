# introduction

This directory contains a collection of template files to build a REST API Forjj plugin.

## REST API Forjj plugin template in GO

When you want to create a Forjj plugin, you can build you own Container
with the service develop in any kind of language.
As Forjj is developped in GO, it was a good opportunity to build Forjj
plugins in GO as well.

If we want to re-use and simplify Forjj plugin development, it became
normal to write plugins from a generic GO code so I introduced the GO
template.

So, this current path contains template files to build/update any Forjj
plugin you want.

If you want to contribute to the template, have a look in
[Developping Forjj template](#Developping_Forjj_template) section.

# How to create a new Forjj plugin in GO?

It is quite easy:

1. Set your GO environment and create your Forjj GO plugin path.
 Assuming GO source code is in ~/go

 ```bash
 $ export GOPATH=~/go
 ```

 Fork and clone forjj-contribs create a 2 level subdirectory.
 forjj-contribs/<pluginType>/<plugin name>

 ```bash
 $ mkdir -p $GOPATH/src/github.hpe.com/christophe-larsonneur # <your Github login>
 $ cd $GOPATH/src/github.hpe.com/christophe-larsonneur # <your Github login>
 $ git clone https://github.hpe.com/christophe-larsonneur/forjj-contribs # <your Github login>
 $ cd forjj-contribs
 $ mkdir -p test/myplugin # <pluginType>/<plugin name>
 ```

2. Create the plugin.go file with 3 comments

```bash
$ vim plugin.go
```

With:

```golang
package main

//go:generate go get github.com/forj-oss/goforjj gopkg.in/yaml.v2
//go:generate go build -o $GOPATH/bin/forjj-genapp github.com/forj-oss/goforjj/genapp
//go:generate forjj-genapp <pluginName>.yaml $GOPATH/src/github.com/forj-oss/goforjj/genapp
```

3. Generate your plugin code

```bash
$ go generate
```

Now you could see a collection of files generated.

## Developping your Forjj plugin from this template

The plugin is generated from a yaml file which defines data exposed to
Forjj cli and handle by the plugin to make the core work.

## FORJJ REST API Reference

Forjj and your Forjj plugin communicate through 3 core actions:

__create action__
POST <url>/create?

payload is in json

__update action__
POST <url>/update?

payload is in json

__maintain action__
POST <url>/maintain?

payload is in json

The json Payload is really simple and formatted as follow:

```json
{"Forj": {
   "<attr>": "<value>",
   ...
   },
 "Objects": {
   "<Object name>": {
     "<Instance Name>": {
       "<Action Name>": {
         "<attr>": "<value>",
         ...
         }
       }
     }
   }
}
```

The content sent by forjj are defined by your plugin yaml file
(<pluginName>.yaml).
The generated code has been generated to conform to this yaml definition
as well.

So, if you change your plugin definition, you need to execute a
`go generate` to update your plugin structure.

**NOTE**: Executing `go generate` several time will update only
`generated-yaml-structs.go`.

**NOTE**: If you removed a generated files, `go generate` will re-create it.

## FORJJ Plugin REST API answer

TBD

# Developing Forjj template

The collection of files used to build the plugin is defined in ../models.go

Any files written here are buildable and runnable but contains only the
code code of the plugin.

So, you can do a `go build` here to verify that your code compiles.

To make it work with `genapp`, some predefined word has been created and
will be replaced by `genapp` to be interpreted by go template to generate
the final plugin source code.

## Tags identified by genapp and replaced for GO template:

- Any lines starting with `.*// __MYPLUGIN: ?` will be replaced by nothing
  Ex:
```golang
type CreateArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $Opts := .Yaml.Objects }}\
	repo map[string]map[string]string // __MYPLUGIN:     {{ go_vars $Objectname}} map[string]map[string]string `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{ end }}\
}
```

will be replaced by
```golang
type CreateArgReq struct {
{{ range $Objectname, $Opts := .Yaml.Objects }}\
    {{ go_vars $Objectname}} map[string]map[string]string `json:"{{$Objectname}}"` // Object details
{{ end }}\
}
```

  So strings before `// __MYPLUGIN: ` are sample code for the local go build
  And strings after `// __MYPLUGIN: ` are go template code

- `__MYPLUGIN__` will be replaced by `{{ go_vars .Yaml.Name }}`
  `go_vars` ensure Yaml.Name is a will be valid GO variable Name.
  It calls [strings.Title()](https://golang.org/pkg/strings/#ToTitle) and remove `-`.
  So a `test-data` will be changed to `TestData`

  Ex:
```golang
func (p *__MYPLUGIN__Plugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
```

will be replaced by
```golang
func (p *{{ go_vars .Yaml.Name }}Plugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
```

- `__MYPLUGINNAME__` will be replaced by `{{ .Yaml.Name }}`
  Usually used in strings. This is not necessarily required for go build
  to work. So, use it if it makes the code more readable.
Ex:
```golang
log.Print("Checking __MYPLUGINNAME__ source code path existence.")
```

will be replaced by
```golang
log.Print("Checking {{ .Yaml.Name }} source code path existence.")
```

- `__MYPLUGIN_UNDERSCORED__` will be replaced by `{{ go_vars_underscored .Yaml.Name }}`
  `go_vars_underscored` function replace any `-` to `_`

Forjj Team
