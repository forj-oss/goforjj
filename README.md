# GO Forjj

## Introduction

This repo contains several golang packages and a golang generator to create a FORJJ plugin and implements the FORJJ plugin protocol in golang.

## What's this repo provides?

2 things:

- go-forjj package
- a GO generate code from yaml file to build/maintain FORJJ plugin.

### go-forjj package

This code provides several generic functions and the Plugin CLI default structure to create your plugin in GO.

list to TDB. Probably built from `godoc`
For now, I suggest you to read the [main goforjj source file] (forjj-plugin-app.go)

## Why creating your forjj plugin

Why will I need to create a forjj plugin?
Well, usually, you won't need it.

But if you want to add an application that is not covered by current FORJJ plugins, then you will certainly need to write your own plugin to manage the new application (create/configure and maintain it) and share it.

You have 2 possibilities:

1. create your own plugin in any language that must respect forjj plugin protocol. It can be a REST API or a simple script which return a json data and can be started from docker or natively where forjj is running.
2. Use `goforjj` to create your plugin in golang.

golang is an awesome language that impressed me when developping forjj.
To help you be focus in your core task (create an app returning json or create a FORJJ REST API ) I built some great fast and powerful golang code to create your first forjj plugin in minutes.

This code will implement the FORJJ plugin protocol.

# Create your FORJJ plugin with GO

1. Create an empty repository. The recommendation is to name it `forjj-<pluginName>`. But this is not a requirement. 
    The core requirement is that your repository name is your plugin name.

    ```bash
    git init forjj-myplugin
    cd forjj-myplugin
    printf "# forjj-myplugin\n\n## Introduction\n\nThis is my first forjj plugin\n" > README.md
    git add README.md
    git commit -m "Initial commit"
    ```

2. If you want your plugin to be part of Forjj ecosystem, the repository needs to be hosted in forj-oss organization.
    To create it, clone git@github.com:forj-oss/forj-oss-infra, and add your repository in the `deployments/production/Forjfile` (`repositories`) and a repository team (`groups/<forjj-<pluginName>/members[]`)
    This team will be owner of the repo.

    You have example in [forj-oss-infra](https://github.com/forj-oss/forj-oss-infra/blob/master/deployments/production/Forjfile)

    To attach your initial git repo to your new forj-oss/forjj-myplugin, do the following:

    ```bash
    git remote add origin git@github.com:<yourName>/forjj-myplugin
    git remote add upstream git@github.com:forjj-oss/forjj-myplugin
    ```

    As soon as forj-oss has created your repo, we should be able to 
3. Optionnally, with GO, a build-env repo helps to build forj-oss projects with your code/build specification. We assume, you are using bash as shell.
  
    Do the following:

    ```bash
    git clone git@github.com:forj-oss/build-env ../build-env
    ../build-env/configure-build-env.sh forjj-myplugin go
    ```

    Load the new build-env (bash script)

    ```bash
    source build-env.sh
    ```
    And create your GO build container

    ```bash
    create-go-build-env.sh
    ```

    As soon as your build-env is loaded, you have direct access via docker of `go`, `glide` and `inenv`

    When you want to unload your build environment, call `build-env-unset`

    When you want to load it again, call `build-env` (alias added your `.bashrc`)

    If you do not take care, the GO binary you will produce will have some dependency to the build environment libraries.
    To avoid that, edit `build-env.sh` and add `export CGO_ENABLED=0` and edit `build-unset.sh` and add `unset CGO_ENABLED`


4.  You need the genapp binary in the /bin of your go tree 
    genapp is part of goforjj project, currently you need to build it from the source.

    [Build genapp instructions](build_genapp.md)
   
5. Move to your plugin directory and create a `plugin.go` with the following:

    ```go
    package main

    //go:generate go build -o /go/bin/genapp forjj-myplugin/vendor/github.com/forj-oss/goforjj/genapp
    //go:generate /go/bin/genapp myplugin.yaml vendor/github.com/forj-oss/goforjj/genapp

    ```

6. Then do :

    ```bash
    glide init           # Initialize of glide.yaml and glide.lock
    glide cc 		 # Optional, Clear glide cache if you have already download a template 
    glide get github.com/forj-oss/goforjj         # Get the goforjj project in te tree to access the template
    ```

    ```bash
    go generate           # To generate the flags management code.
    ```
    Depending on your plugin definition (`/runtime/service_type`), `go generate` will create several files:

    ```bash
    glide update           # Update dependencies
    ```

    the initial generated code is functionnal! So, you can do:

    ```bash
    go build           # Build plugin
    ```

  Now you should have a hello-world plugin 

- `REST API` : Default.

If you have added some extra commands, a `<command>.go` will be also created the first time with initial code. So, you will just need to edit it and add your specific command code.

**NOTE**: A lot of features and flags values are managed by the `goforjj` package, please read the `goforjj` package [documentation above] (#go-forjj-package).

## Writing your `<plugin>.yaml` file

This file is the plugin declarative part of your plugin.

It defines what your plugin exposes and how forjj do task with your plugin.

If you run a FORJJ plugin from the genapp templates (go generate), this file is also used to generate/maintain some GO source files.

The yaml file must be structured as follow:

```yaml
---
plugin: <PluginName>       // Required
version: <PluginVersion>   // Optional
description: <description> // Optional
created_flag_file: "{{ .InstanceName }}/{{.Name}}.yaml" # Optional. By default "{{ .InstanceName }}/{{.Name}}.yaml". Usually is the <plugin>.yaml source file stored in the infra source repository. Ex: github/github.yaml
runtime:   // Define how the plugin is started
  docker:  // This is the default and only way to get it working today.
    image: "string" // Docker Image name to use
  service_type: "" // "REST API" or "shell"
  service:
    socket: "string" // Optional: Name of the socket to use. By default, it creates <PluginName>.sock
    parameters: [ "string", ... ] // Optional: Collection of parameters to start the plugin service.
actions: // Obsolete. Is replaced by task_flags and objects.
  common: // Flags defined for each tasks
    <FlagName>: // Please note the spaces shift for next lines. It have to be indented from the <FlagName> string position.
      help: "<Help string>"  // If missing, no help is displayed.
      required: <false/true> // If missing, required is false.
      hidden: <false/true>   // if missing, hidden is false.
      group: "string"        // Used in GO genapp templates to regroup several fields in a struct.
      short: '<caracter>'    // single caracter for short option. if Missing no short option set.
      secure: <false/true>   // false by default. If the flag is given, forjj will save in a `forjj-creds.yml` like file in your workspace. The plugin should not save it anywhere.
      [...]                    // Collection of FlagName...
  create:   // Flags defined for 'create' task. Same syntax as found in 'common' section.
      ...
  update:   // Flags defined for 'update infra' infra task. Same syntax as found in 'common' section.
      ...
  maintain: // Flags defined for 'maintain' task. Same syntax as found in 'common' section.
      ...
task_flags:
  common: // Flags defined for each tasks
    <FlagName>: // Please note the spaces shift for next lines. It have to be indented from the <FlagName> string position.
      help: "<Help string>"  // If missing, no help is displayed.
      required: <false/true> // If missing, required is false.
      hidden: <false/true>   // if missing, hidden is false.
      short: '<caracter>'    // single caracter for short option. if Missing no short option set.
      secure: <false/true>   // false by default. If the flag is given, forjj will save in a `forjj-creds.yml` like file in your workspace. The plugin should not save it anywhere.
      [...]                    // Collection of FlagName...

  add:     // Flags defined for 'add' or others tasks. Same syntax as found in 'common' section.
  remove:
  rename:
  list:
  maintain:
objects:          // Collection of options that the plugin expose to forjj.
  <object_name>:
    actions :     // By default: [ "add", "update", "remove", "rename", "list"]. If you limit to few actions, set it here.
    flags:        // Collection of flags for the object to be managed by forjj.
      <flagName>:
        only-for-actions: ["string", ...] // If missing, use the actions defined at object level.
        help: "<Help string>"             // If missing, no help is displayed.
        required: <false/true>            // If missing, required is false.
        hidden: <false/true>              // if missing, hidden is false.
        short: '<caracter>'               // single caracter for short option. if Missing no short option set.
```

### Short examples

In order to start defining your own plugin declaration, here is few examples to highlight how to set it.

First of all, here is a minimal version:

```yaml
plugin: my-plugin
runtime:
  docker:
    image: "myimage" // Docker Image name to use
```

But to be honest, with this yaml file, it will work, for sure, because you ask forjj to do ... nothing.

So, you will need to enhance it with at least, one flag

```yaml
plugin: my-plugin
runtime:
  docker:
    image: "myimage" // Docker Image name to use
actions:
  add:
    my-flag:
      help: Help about my flag
      required : true
```

This one expose 1 flag that forjj will list as possible and REQUIRED when we do a create.

So, you could ask forjj to provide this flag as follow:

```bash
forjj create ~/tmp/my-workspace --apps mycateg:myplugin --myplugin-my-flag flag-value
forjj create ~/tmp/my-workspace --apps mycateg:myplugin:myinstance --myinstance-my-flag flag-value
forjj add repo mycateg:myplugin --my-flag flag-value
forjj add repo mycateg:myplugin:myinstance --my-flag flag-value
```

In this case, forjj will create a workspace, add your plugin in the list of managed applications and provide the `flag-value` to the flag name `my-first-flag`

You can get some forjj internal variables with `forjj-<data>`
Ex:


```yaml
plugin: my-plugin
runtime:
  docker:
    image: "myimage" // Docker Image name to use
actions:
  add:
    my-flag:
      help: Help about my flag
      required : true
    forjj-organization:
```

In this example, forjj will provide a forjj-organization flag to your plugin with the organization name created by forjj.

Under `actions`, you can set `common`, `create`, `update` and `maintain`

Except `common`, those entries are main forjj actions for forjj infra object type.

`common` helps you to declare same flags for create and update.

Ex: --my-first-flag usable for `create` and `update infra` tasks

```yaml
plugin: my-plugin
runtime:
  docker:
    image: "myimage" // Docker Image name to use
actions:
  common:
    my-common-flag:
      help: Help about my flag
      required : true
```

You can also expose some objects that your plugin will manage and that you want forjj to expose to end users.

Ex:

```yaml
plugin: my-plugin
runtime:
  docker:
    image: "myimage" // Docker Image name to use
objects:
  user:
    flags:
      my-first-flag:
        help: Help about my flag
        required : true
```

In this case, you can `add`, `update`, `remove`, `rename`, `list` one or more user.

Ex:
```bash
forjj add user --my-first-flag flag-value
forjj remove user --my-first-flag flag-value
forjj list user
forjj create ~/tmp/my-workspace --apps categ:my-plugin
```

You can limit some flags to less actions:

```yaml
plugin: my-plugin
runtime:
  docker:
    image: "myimage" // Docker Image name to use
objects:
  user:
    flags:
      my-first-flag:
        only-for-actions: ["add", "remove"]
        help: Help about my flag
        required : true
```
