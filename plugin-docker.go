package goforjj

import (
    "fmt"
    "os"
    "time"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
    "strings"
    "regexp"
)

type DockerService struct {
    Volumes map[string]byte
    Env map[string]byte
}

// Define how to start
func (p *PluginDef) docker_start_service(instance_name string) (is_docker bool, err error) {
    if p.Yaml.Runtime.Docker.Image == "" {
        return false, nil
    }
    is_docker = true
    p.docker.name = instance_name
    gotrace.Trace("Starting it as docker container '%s'", p.docker.name)

    // intialize
    p.docker.init()

    // mode daemon
    p.docker.opts = append(p.docker.opts, "-d")

    // Source path
    if _, err := os.Stat(p.Source_path) ; err != nil {
        os.MkdirAll(p.Source_path, 0755)
    }
    p.SourceMount = "/src/"
    p.docker.add_volume(p.Source_path + ":" + p.SourceMount)

    // Workspace path
    if p.Workspace_path != "" {
        p.WorkspaceMount = "/workspace/"
        p.docker.add_volume(p.Workspace_path + ":" + p.WorkspaceMount)
    }

    // Do we have a socket to prepare?
    if p.Yaml.Runtime.Service.Port == 0 && p.cmd.socket_path != "" {
        if err = p.socket_prepare(); err != nil {
            return
        }
        p.docker.socket_path = "/tmp/forjj-socks"
        p.docker.add_volume(p.cmd.socket_path+":"+p.docker.socket_path)
        p.docker.opts = append(p.docker.opts, "-u", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
    } else {
        err = fmt.Errorf("Forjj connect to remote url - Not yet implemented\n")
        return
    }

    if p.Yaml.Runtime.Docker.Volumes != nil {
        for _, v := range p.Yaml.Runtime.Docker.Volumes {
            p.docker.add_volume(v)
        }
    }

    if p.Yaml.Runtime.Docker.Env != nil {
        for _, v := range p.Yaml.Runtime.Docker.Env {
            p.docker.add_env(v)
        }
    }

    if p.Yaml.Runtime.Docker.Dood {
        if p.dockerBin == "" {
            err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path", p.docker.name)
            return
        }
        p.docker.add_volume("/var/lib/docker/docker.sock:/var/lib/docker/docker.sock")
        p.docker.add_volume(p.dockerBin + ":/bin/docker")
        // TODO: download bin version of docker and mount it, or even communicate with the API directly in the plugin container (go: https://github.com/docker/engine-api)

    }

    // Check if the container exists and is started.
    // TODO: Be able to interpret some variables written in the <plugin>.yaml and interpreted here to start the daemon correctly.
    // Ex: all p.cmd_data .* in a golang template would give {{ .socket_path }}, etc...
    if _, err = p.docker_container_restart(p.cmd.command, p.Yaml.Runtime.Service.Parameters); err != nil {
        return
    }

    err = p.check_service_ready()

    return
}

// Regularly testing the service response. fails after a timeout.
func (p *PluginDef) check_service_ready() (err error) {
    gotrace.Trace("Checking service status...")
    for i := 1; i < 30; i++ {
        time.Sleep(time.Second)

        out := ""
        if out, err = docker_container_status(p.docker.name) ; err != nil {
            return
        }

        if strings.Trim(out, " \n") != "running" {
            out, err = docker_container_logs(p.docker.name)
            if  err == nil {
                out = fmt.Sprintf("docker logs:\n---\n%s---\n",out)
            } else {
                out = fmt.Sprintf("%s\n", err)
            }
            docker_container_remove(p.docker.name)
            err = fmt.Errorf("%sContainer '%s' has stopped unexpectedely.", out, p.Yaml.Name)
            return
        }

        if p.CheckServiceUp() {
            gotrace.Trace("Service is UP.")
            p.service_booted = true
            return
        }
    }
    err = fmt.Errorf("Plugin Service '%s' not started successfully as docker container '%s'. check docker logs\n", p.Yaml.Name, p.docker.name)
    return
}

func (d *docker_container) init() {
    d.opts = make([]string, 0, 5)
    d.volumes = make(map[string]byte)
    d.envs = make(map[string]byte)
}

func (d *docker_container) add_volume(volume string) {
    if ok, _ := regexp.MatchString("^.*(:.*(:(ro|rw))?)?$", volume ) ; ok {
        d.volumes[volume] = 'v'
    }
}

func (d *docker_container) add_env(env string) {
    if ok, _ := regexp.MatchString("^.*=.*$", env ) ; ok {
        d.envs[env] = 'e'
    }
}

func (d *docker_container) complete_opts_with(val ...map[string]byte) {
    // Calculate the expected array size
    tot := len(d.opts)

    for _, v := range val {
        tot += len(v)*2
    }

    // allocate
    r := make([]string, 0, tot)

    // append
    r = append(r, d.opts...)
    for _, v := range val {
        for k, o := range v {
            r = append(r, "-" + string(o))
            r = append(r, k)
        }
    }
    d.opts = r
}
