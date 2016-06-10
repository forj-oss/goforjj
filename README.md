# Introduction

This repo contains a GO module to read and expose a FORJJ plugin flags from a yaml file to the output and a `flags` option of the plugin tool built on GO.

# How to use it?

In your main go file, you should add:

    ```go
    import "github.hpe.com/christophe-larsonneur/go-forjj"

    //go:generate go run $GOPATH/github.hpe.com/christophe-larsonneur/go-forjj/cmd/genflags/main.go <PluginName>.yaml

    ```

Then do :

    ```bash
    go get                # Download dependencies, and by the way `go-forjj`
    vim <PluginName>.yaml # To define the plugin flags input.
    go generate .         # To generate the flags management code.
    vim main.go           # To call the `<PluginName>App.New(os.Args[1:])`
    ```

# What's this repo provides?

2 things:

- go-forjj module
- a GO generate code from yaml file

##  go-forjj module

This code provides several generic functions to create your plugin in GO.

list to TDB

## GO generate code

As soon as you add the `go:generate` comment, you need to generate a go source file if you updated your `<YamlFileName>.yaml` with:

    ```bash
    go generate .
    ```

This will produce a `<PluginName>_generated.go` file which will contains the plugin flag code using kingpin module.

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

   
