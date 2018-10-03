# Forjj plugin documentation

`forjj` thanks to `Forjfiles` will `create`/`update` and `maintain` the software factory that you defined.
Your `Forjfile` requires to define at least one or more application (`applications/<ApplicationName>`) which will compose your Software factory.

`forjj` can't manage them if there is no *forjj plugin* attached to. The forjj plugin is also known as driver.

This documentation explains how forjj detect, download, use them and what a plugin typically do in such action cases.

## How forjj identify a plugin to manage an application

An application is defined in your Forjfile as follow:

```yaml
applications: # List all applications
  myAppName: # required. This is the application name
    driver: mydriver # Optional. If not set, forjj will take the application name as driver name
    type: upstream # Required
    [...] # collections of application parameters. Some can be required, with default value or optional.
```

`applications/<appName>/driver` (or `<appName>` if not set) is the name of the forjj plugin that forjj will use to manage the application through `create`, `update` or `maintain` action.

As soon as `forjj` has loaded all Forjfiles, it build a complete list of drivers to load and use.

To load them, it will search for drivers definition to load as follow:

It will search them from a list of urls which is described by `forjj workspace` `contrib-repo-path`.
By default, forjj will search in `https://github.com/forj-oss/<repo>/raw/master`. `<repo>` is replaced by `<driverName>` or `forjj-<driverName>`. So Forjj will check those path to find and download the `<pluginName>.yaml`

For example, for `github` plugin, forjj will search for :

- `https://github.com/forj-oss/github/raw/master`/`github.yaml`
- `https://github.com/forj-oss/forjj-github/raw/master`/`github.yaml`

First found, first loaded. Forjj will exit if one driver definition is missing.

**NOTE** Please note that forjj run a `forjj validate` to verify the Forjfile setup with plugins definition. Forjj can report you some errors if you set some data which has no effect as not recognized by forjj and loaded plugins.

## How forjj determine how to download and start the plugin

When forjj has loaded a plugin definition, it search for the runtime section which describes where is the plugin image and how to start it.

Typical `<pluginName>.yaml`:

```yaml
[...]
runtime: # runtime section to describe how to download and start the plugin.
  service_type: "REST API" # Required. Support only 'REST API'.
  docker:
    image: "forjdevops/forjj-myplugin" # Required to run the plugin from docker. Describe the docker image to pull. ie: docker pull <image>. It respect the docker pull image string format.
    # dood: # False by default. Set true if your plugin requires docker in it. forjj will configure the plugin container to run docker from it in a Docker Out Of Docker mode. Explained later.
    # volume: # Optional. Array of volume to add (docker -v)
    # env: # Optional. Array of environment to set (docker -e)
    # user: # Optional. User to start the process (docker -u)
  service:
    # socket: "myplugin.sock" # by default, `<pluginName>.sock` if port is empty
    # port: # Not yet implemented. Should be used to set the docker tcp model.
    # command: # string. Command to run to start the plugin.
    parameters: [ "service", "start", "--templates", "/templates"] # Array of parameters to provide to the plugin to start.
[...]
```

`docker` section describes how to find and pull the plugin image from docker and will use docker to start it.

The runtime/docker/image is required as the image will be pulled by docker automatically by forjj thanks to this field.

**NOTE**: if you are running your plugin in debug mode, you can start your plugin by hand, outside docker. Use `forjj --run-plugin-debugger=Plugin1[,plugin2[,...]]`

## Plugin docker image design

`Forjj` starts plugins from docker. If those plugins requires to run some docker task, the plugin container will need to be configured with DooD (Docker out of Docker) or DinD (Docker in Docker)

You can read some documentation about those 2 differents technics [here](http://blog.teracy.com/2017/09/11/how-to-use-docker-in-docker-dind-and-docker-outside-of-docker-dood-for-local-ci-testing/).

DooD is one recommended way to run docker commands in docker container. So `Forjj` has implemented some features to manage those container in DooD mode, automatically. It will be explained in [next section](#plugin-dood-docker-image-design).

If your plugin do not use `Docker DooD`, then `forjj` will create a container with a minimum set of parameters.

### Minimal docker run parameters used by forjj to start a plugin

`forjj` will run the plugin with `docker run`, to create the container or with `docker start` to restart the container.

The container will have the following mount:

- /tmp/forj_socks : where the plugin will create the socket, so that forjj can communicate with the plugin
- /workspace      : mount of the infra/.forj-workspace directory
- /src            : mount of plugin source code stored in your infra repository.
- /deploy         : mount of deployments parent directory, containing all deployment reposistories.

if forjj is called with `FORJJ_SOURCE_BASE`, mounts are different:

- /tmp/forj_socks    : where the plugin will create the socket, so that forjj can communicate with the plugin
- $FORJJ_SOURCE_BASE : unique and base directory which contains src, deploy and workspace. Jenkins service uses that variable to mount the jenkins_home.
                       Forjj will take this global path in account to determine where sources, deploys, and workspace are located in the plugin container

The container will have the following environment:

- http_proxy/https_proxy/no_proxy : if set in your workstation.
- PROXY                           : a docker run shell helper for proxies. To use it in shell, you will need to run `eval "docker run $PROXY"`
  - `-e http_proxy`               : http proxy setup
  - `-e https_proxy`              : https proxy setup
  - `-e no_proxy`                 : no proxy setup
- LOGNAME                         : Username used to run forjj.

The container will be started with :

- the user/group ID used to start forjj. (docker -u UID:GID)
- default directory (pwd) to /src (or source relative path if FORJJ_SOURCE_BASE is set). SELF_SRC, SELF_DEPLOY and SELF_WORKSPACE are set to the container path.
- command defined by `runtime/service/command`. If not set, you will need to define it with CMD or ENTRYPOINT in your Dockerfile.
- command parameters defined by `runtime/service/command`

For such simple plugin, the plugin image needs to store the plugin binary or tool which will create the socket and listen for tasks from `forjj`

### plugin Dockerfile

As this kind of plugin are quite simple, your Dockerfile should be simple as well, depending on the tool that you want to run to listen to the socket channel with `forjj`

Initial forjj plugins (forjj-jenkins / forjj-github) were written in GO. 
They were generated from some template to create a vanilla plugin, doing nothing but running. So you can just start writing the plugin logic to manage your application. code 
You can find [some explanation how to create a basic one](README.md#create-your-forjj-plugin-with-go) which create a simple static binary.

```Forjfile
FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates --fresh && \
    rm -f /var/cache/apk/*tar.gz && \
    adduser devops devops -D

COPY forjj-aPlugin /bin/aPlugin

USER devops

ENTRYPOINT ["/bin/aPlugin"]

CMD ["--help"]
```

In this example, we are building the GO plugin separately and copy the binary to the docker image with `COPY`

## Plugin DooD docker image design

This image is much more complex as it introduce Docker out of Docker (DooD) concept in the plugin container.

To activate it, the plugin must set `runtime/docker/dood` to true.
With this option, your plugin will be able to run `docker` command and will use the forjj host `docker cache`.

Then forjj will configure the plugin container as follow:

### Docker run parameters used by forjj to start a plugin DooD

`forjj` will run the plugin with `docker run` to create the container or with `docker start` to restart the container.

The container will have the following mount:

- /tmp/forj_socks      : where the plugin will create the socket, so that forjj can communicate with the plugin
- /workspace           : mount of the infra/.forj-workspace directory
- /src                 : mount of plugin source code stored in your infra repository.
- /deploy              : mount of deployments parent directory, containing all deployment reposistories.
- /var/run/docker.sock : DooD - Host Docker socket
- /usr/bin/docker      : DooD - static binary as described by `forjj workspace` `docker-bin-path`

The container will have the following environment:

- http_proxy/https_proxy/no_proxy : if set in your workstation.
- PROXY                           : a docker run shell helper for proxies. 
                                    To use it in shell, you will need to run `eval "docker run $PROXY"`.
  - `-e http_proxy`               : http proxy setup
  - `-e https_proxy`              : https proxy setup
  - `-e no_proxy`                 : no proxy setup
- LOGNAME                         : Current user name used to run forjj.
- UID                             : DooD - Current user ID which has started forjj.
- GID                             : DooD - Current user group ID which has started forjj.
- DOCKER_DOOD_GROUP               : DooD - Group ID of the docker socket file. We assume name to be `docker`.
- DOOD_SOURCE                     : DooD - docker run shell helper for sources mount. if forjj uses `FORJJ_SOURCE_BASE`, this shell string will just have a single `-v $FORJJ_SOURCE_BASE:$FORJJ_SOURCE_BASE`
                                    To use it in shell, you will need to run `eval "docker run $DOOD_SOURCE"`.
- DOCKER_DOOD                     : DooD - String of docker run options to mount and set environment. Used to run a DooD container from a DooD container.
                                    To use it in shell, you will need to run `eval "docker run $DOCKER_DOOD_BECOME"`.
  The list of options are:
  - `-v <hostDockerSocket>:/var/run/docker.sock`
  - `-v <hostDockerBinPath>:/usr/bin/docker`
  - `-e DOOD_SRC=<hostInfraPluginSource>`
  - `-e DOOD_DEPLOY=<hostPluginSource>`
  - `-e DOCKER_DOOD_GROUP=<hostDockerGroup>`
- DOCKER_DOOD_BECOME              : DooD - String of docker run option to become root and set environment variable UID/GID/DOCKER_DOOD_GROUP.
                                    To use it in shell, you will need to run `eval "docker run $DOCKER_DOOD_BECOME"`.
  In details:
  - `-u root:root`
  - `-e UID=<hostCurrentUserUID`
  - `-e GID=<hostCurrentUserGID`

**NOTE**: UID/GID can be set outside DooD Context, if the container started as root needs to become a user with a different UID/GID.

The container will be started with :

- the user/group ID used to start forjj. (docker -u UID:GID)
- default directory (pwd) to /src
- command defined by `runtime/service/command`. If not set, you will need to define it with CMD or ENTRYPOINT in your Dockerfile.
- command parameters defined by `runtime/service/command`

### plugin Dockerfile

In this DooD Context, The plugin image must take care of the DooD environments variables given 
(UID, GID, DOCKER_DOOD_GROUP, DOCKER_DOOD & DOCKER_DOOD_BECOME)

When the forjj plugin container start, if:

- UID & GID are set:

  - the plugin process will need to be started with those UID & GID
  - if needed, the current container user must be updated with those UID/GID
- DOCKER_DOOD_GROUP is set:
  - a docker group must exist or update with the id given in this variable.
- DOCKER_DOOD, DOCKER_DOOD_BECOME are set or not:
  - nothing to do, but if that container needs to create a container in DooD mode, those variables can be used as is to the docker run command. Ex: `eval "docker run $DOCKER_DOOD $DOCKER_DOOD_BECOME [...]"`

To simplify your entrypoint script, you can use [`docker-lu`](https://github.com/forj-oss/docker-lu) This program update passwd and group file with needed values

Example of a Dockerfile for DooD container:

```Dockerfile
FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

ADD https://github.com/forj-oss/docker-lu/releases/download/0.1/docker-lu /usr/local/bin/docker-lu

RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates --fresh && \
    rm -f /var/cache/apk/*tar.gz && \
    adduser devops devops -D && \
    chmod 755 /usr/local/bin/docker-lu

COPY forjj-aPlugin /bin/aPlugin
COPY dockerfiles/entrypoint.sh /usr/local/bin/entrypoint.sh


USER devops

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

CMD ["--help"]
```

entrypoint.sh:

```sh
#!/bin/sh

docker-lu -u devops $UID -g devops $GID -G docker $DOCKER_GROUP
exec /bin/su devops -c "/bin/aPlugin $@"
```

With UID/GID/DOCKER_DOOD_GROUP, you can use `sed` and `groupadd`/`addgroup` depending on the linux release used to create/update properly. 
But you can use instead which do this in a single line more securily: [`docker-lu`](https://github.com/forj-oss/docker-lu) was written for that perspective.

## minimal plugin yaml definition

When `forjj` works with plugins, it sends a REST API post request, to `create`/`update` or `maintain` applications by those plugins.

The plugin is started from docker and files are made available through docker mount mechanism and it will need to read files from a source directory and write to a destination directory.

Each request payload is dynamically built by forjj from the plugin file definition with data from Forjfiles, forjj internal and forjj secrets.
At least, the plugin requires forjj to identify where are plugin source files and where to write files (generate, ...)

Forjj internal data:

- forjj-organization:    Forjj oranization name
- forjj-instance-name:   Name of the instance of the current application

- forjj-infra:           Name of the Infra repository to use
- forjj-infra-upstream:  Address of the infra repository upstream

- forjj-deploy-mount:    Container path to the deployment repository
- forjj-source-mount:    Container path to the source directory.
- forjj-workspace-mount: Container path to the forjj workspace.

- forjj-username:        User name running forjj (from $LOGNAME)

Those paths are request by the plugin, by defining them in the `<plugin>.yaml` under section `task_flags/<common|create|update|maintain>/`

### Typical create Request plugin actions

when forjj request a plugin to `create` the application, the plugin needs to:

- create initial source files in the `forjj-source-mount`. If it already exist, it must ask the user to use `update` instead.
- create the deployment code from the initial source files. Generated files will needs to be made in `forjj-deploy-mount`

### Typical update Request plugin actions

when forjj request a plugin to `update` the application, the plugin needs to:

- read source files from the `forjj-source-mount`. If not found, it must request the user to create it first.
- update the deployment code from the source files. Generated files will needs to be made in `forjj-deploy-mount`

### Typical maintain Request plugin actions

when forjj request a plugin to `maintain` the application, the plugin needs to:

- read and execute deployement files stored in `forjj-deploy-mount`

### Using docker in the plugin - Docker Mount context

As the plugin is configured to support DooD thanks to `forjj`, it can start docker or even shell script to run docker commands.

When the plugin has been started, `forjj` has already configured the plugin container with required mounts and environment variables.
But when the plugin needs to run docker containers within the plugin container, mounts cannot use local container path to set a mount,
because the host reference is the host itself and the container.

Here is an example: 
If the plugin container has been started with `-v /var/container/home:/home`, `/home`, an ls in the plugin container will shows files/directories stored in `/var/container/home` of the host.

if you start a container from this container, with `-v /home:/home`, the container created will mount the HOST `/home` to `/home` in the new container.
The main confusion is to remind where docker daemon run. 
If dockerd is running in the host (DooD), the parent mount path reference is the HOST mount path. Not the container mount path.
if dockerd is running in the container (DinD), the parent mount path reference is the CONTAINER mount path. Not the host mount path

To facilitate running container in the same context as the plugin, `forjj` delivers few additionnal environment variables which was already listed above in section [Docker run parameters used by forjj to start a plugin DooD](#Docker-run-parameters-used-by-forjj-to-start-a-plugin-DooD)

- DOCKER_DOOD        : docker run options to configure the new container is DooD
- DOCKER_DOOD_BECOME : docker run options to run as root and become a different user (UID/GID)
- PROXY              : docker run options to add host proxy setup
- DOOD_SOURCE        : docker run options to mount source/deploy/workspace path

If the plugin container needs to run a container with one or more `docker run` shell helper, you have to run it in an `eval`

ex: 

    eval "docker run $DOOD_SOURCE"
Enjoy

Forj team