package goforjj

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
    "net"
    "net/url"
    "os"
    "path"
    "time"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
    "strings"
)

const defaultTimeout = 32 * time.Second

// This function start the service as daemon and register it
// If the service is already started, just use it.
func (p *PluginDef) PluginStartService() error {
    if !p.service { // Nothing to start
        return nil
    }
    gotrace.Trace("Starting plugin service...")
    // Is it a docker service?
    if p.Yaml.Runtime.Image != "" {
        p.docker.name = p.Yaml.Name
        gotrace.Trace("Starting it as docker container '%s'", p.docker.name)

        // mode daemon
        p.docker.opts = []string{ "-d" }
        // Source path
        p.SourceMount = "/src/"
        p.docker.opts = append(p.docker.opts, "-v", p.Source_path + ":" + p.SourceMount,)

        // Workspace path
        if p.Workspace_path != "" {
            p.WorkspaceMount = "/workspace/"
            p.docker.opts = []string{"-v", p.Workspace_path + ":" + p.WorkspaceMount}
        }
        // Do we have a socket to prepare?
        if p.Yaml.Runtime.Service.Port == 0 && p.cmd.socket_path != "" {
            if err := p.socket_prepare(); err != nil {
                return err
            }
            p.docker.socket_path = "/tmp/forjj-socks"
            p.docker.opts = append(p.docker.opts, "-v", p.cmd.socket_path+":"+p.docker.socket_path)
            p.docker.opts = append(p.docker.opts, "-u", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
        } else {
            return fmt.Errorf("Forjj connect to remote url - Not yet implemented\n")
        }

        // Check if the container exists and is started.
        // TODO: Be able to interpret some variables written in the <plugin>.yaml and interpreted here to start the daemon correctly.
        // Ex: all p.cmd_data .* in a golang template would give {{ .socket_path }}, etc...
        if _, err := p.docker_container_restart(p.cmd.command, p.Yaml.Runtime.Service.Parameters); err != nil {
            return err
        }

        gotrace.Trace("Checking service status...")
        for i := 1; i < 30; i++ {
            time.Sleep(time.Second)

            var err error
            out := ""

            if out, err = docker_container_status(p.docker.name) ; err != nil {
                return err
            }

            if strings.Trim(out, " \n") != "running" {
                out, err = docker_container_logs(p.docker.name)
                if  err == nil {
                    out = fmt.Sprintf("docker logs:\n---\n%s---\n",out)
                } else {
                    out = fmt.Sprintf("%s\n", err)
                }
                docker_container_remove(p.docker.name)
                return fmt.Errorf("%sContainer '%s' has stopped unexpectedely.", out, p.Yaml.Name)
            }

            if p.CheckServiceUp() {
                gotrace.Trace("Service is UP.")
                p.service_booted = true
                return nil
            }
        }

        return fmt.Errorf("Plugin Service '%s' not started successfully as docker container '%s'. check docker logs\n", p.Yaml.Name, p.docker.name)
    }
    // Run a command which should become a daemon (options or shell fork)

    // Do we have a socket to prepare?
    if p.Yaml.Runtime.Service.Port == 0 && p.cmd.socket_path != "" {
        if err := p.socket_prepare(); err != nil {
            return err
        }
        cmd_args := []string{p.cmd.command}
        cmd_args = append(cmd_args, p.cmd.args...)
        _, err := cmd_run(cmd_args)
        return err
    } else {
        return fmt.Errorf("Forjj connect to remote url - Not yet implemented\n")
    }
}

func (p *PluginDef) CheckServiceUp() bool {
    if p.cmd.socket_file != "" {
        if _, err := os.Stat(path.Join(p.cmd.socket_path, p.cmd.socket_file)); os.IsNotExist(err) {
            return false
        }
    }
    p.url.Path = "ping"
    _, body, err := p.req.Get(p.url.String()).End()
    if err != nil {
        fmt.Printf("Issue to ping the Plugin service '%s'. %s\n", p.Yaml.Name, err)
    }
    return strings.Trim(body, " \n") == "OK"
}

// Create the socket link for http and his path.
func (p *PluginDef) socket_prepare() (err error) {
    p.cmd.socket_file = p.Yaml.Name + ".sock"
    socket := path.Join(p.cmd.socket_path, p.cmd.socket_file)
    p.req = gorequest.New()
    p.req.Transport.Dial = func(_, _ string) (net.Conn, error) {
        return net.DialTimeout("unix", socket, defaultTimeout)
    }
    p.url, err = url.Parse("http://forjj.sock")

    // TODO: Test deeper about the path access.
    _, err = os.Stat(p.cmd.socket_path)
    if err != nil {
        if os.IsNotExist(err) {
            if err := os.Mkdir(p.cmd.socket_path, 0755); err != nil {
                return err
            }
        }
    }
    return nil
}

// To stop the plugin service if the service was started before by goforjj
func (p *PluginDef) PluginStopService() {
    if ! p.service || ! p.service_booted {
        return
    }
    p.url.Path = "quit"
    p.req.Get(p.url.String()).End()

    if p.Yaml.Runtime.Image != "" {
        for i := 0 ; i <= 10 ; i++ {
            time.Sleep(time.Second)
            if out, _ := docker_container_status(p.docker.name); out != "started" {
                return
            }
        }
        if out, _ := docker_container_status(p.docker.name); out == "started" {
            docker_container_stop(p.docker.name)
        }
    }
}
