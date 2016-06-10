# Introduction

This repo contains a GO package to generate, read and expose a FORJJ plugin flags from a yaml file to the output and a `flags` option of the plugin tool built on GO.

# How to use it?

1. Write your FORJJ plugin yaml file as described in [forjj-contribs README](https://github.hpe.com/christophe-larsonneur/forjj-contribs#description-of-yaml)
2. create the `plugin.go` with the following:

    ```go
    import "github.hpe.com/christophe-larsonneur/go-forjj"

    //go:generate go run $GOPATH/github.hpe.com/christophe-larsonneur/go-forjj/cmd/genflags/main.go <PluginName>.yaml

    ```

3. Then do :

    ```bash
    go get                # Download dependencies, and by the way `go-forjj`
    go generate           # To generate the flags management code.
    ```
  The `go generate` will create 4 files:

- `<pluginName>.go` - Plugin core code. Do not update this file. It will be regenerated from your plugin yaml file.
- `<create|update|maintain>.go` - Those files are created the first time. It won't be re-generated. This is where you will have to write your plugin code.

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

   
