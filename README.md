# Introduction

This repo contains several golang packages and a golang generator to create a FORJJ plugin and implements the FORJJ plugin protocol in golang.

# Why creating your forjj plugin

Why will I need to create a forjj plugin?
Because you want forjj to manage an application (create/configure and maintain it) that is currently not supported

You have 2 possibilities:
1. create your own plugin in any language that must respect forjj plugin protocol. It can be a REST API or a simple script which return a json data and can be started from docker or natively where forjj is running.
2. Use `goforjj` to create your plugin in golang.

golang is an awesome language that impressed me when developping forjj.
To help you be focus in your core task (create an app returning json or create a FORJJ REST API ) I built some great fast and powerful golang code to create your first forjj plugin in minutes.

This code will implement the FORJJ plugin protocol.

# Create your FORJJ plugin

1. Optional. Write your FORJJ plugin yaml file as described in [forjj-contribs README](https://github.hpe.com/christophe-larsonneur/forjj-contribs#description-of-yaml)

This file mainly will define a list of flags to provide to the plugin through `forjj` cli

**If you do not create it, go generate will create one for you.**

2. create the `plugin.go` with the following:

    ```go
    package main

    //go:generate go get github.hpe.com/christophe-larsonneur/goforjj gopkg.in/yaml.v2
    //go:generate go build -o $GOPATH/bin/genapp github.hpe.com/christophe-larsonneur/goforjj/genapp
    //go:generate genapp <PluginName>.yaml

    ```
    Replace `<PluginName>.yaml` by your own plugin definition yaml file created at step 1.

3. Then do :

    ```bash
    go generate           # To generate the flags management code.
    ```
  Depending on your plugin definition (`/runtime/service_type`), `go generate` will create several files:

- `shell` :
- `REST API` :

If you have added some extra commands, a `<command>.go` will be also created the first time with initial code. So, you will just need to edit it and add your specific command code.

**NOTE**: the initial generated code is functionnal!!! So, you can do a `go build` and run the binary generated to make a basic try!

So, now you can start your plugin development!!!

**NOTE**: A lot of features and flags values are managed by the `goforjj` package, please read the `goforjj` package [documentation above] (#go-forjj-package).

# What's this repo provides?

2 things:

- go-forjj package
- a GO generate code from yaml file

##  go-forjj package

This code provides several generic functions and the Plugin CLI default structure to create your plugin in GO.

list to TDB. Probably built from `godoc`
For now, I suggest you to read the [main goforjj source file] (forjj-plugin-app.go)

## GO generate code

As soon as you add the `go:generate` comment, you need to generate a go source file if you updated your `<YamlFileName>.yaml` with:

    ```bash
    go generate .
    ```

This will produce a `<PluginName>_generated.go` file which will contains the plugin flag code using kingpin package.

This source file will need to be stored on GIT. `go generate .` is not intended to be executed by any CI. This is the role of the developper to ensure this generated code is correct.

The yaml file must be structured as follow:

    ```yaml
    ---
    plugin: <PluginName>       // Required
    version: <PluginVersion>   // Optional - NOT IMPLEMENTED
    description: <description> // Optional - NOT IMPLEMENTED
    flags:
      common: // Flags defined for each tasks
        - <FlagName>: // Please note the spaces shift for next lines. It have to be indented from the <FlagName> string position.
            help: "<Help string>"  // If missing, no help is displayed.
            required: <false/true> // If missing, required is false.
            hidden: <false/true>   // if missing, hidden is false.
            short: '<caracter>'    // single caracter for short option. if Missing no short option set.
        - [...]                    // Collection of FlagName...
      check:    // Flags defined for 'check' task. Same syntax as found in 'common' section.
        ...
      create:   // Flags defined for 'create' task. Same syntax as found in 'common' section.
        ...
      update:   // Flags defined for 'update' task. Same syntax as found in 'common' section.
        ...
      maintain: // Flags defined for 'maintain' task. Same syntax as found in 'common' section.
        ...
    ```


