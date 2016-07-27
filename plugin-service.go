package goforjj

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
    "net"
    "net/url"
    "os"
    "path"
    "time"
)

const defaultTimeout = 32 * time.Second

// This function start the service as daemon and register it
// If the service is already started, just use it.
func (p *PluginDef) PluginStartService() error {
    if !p.service { // Nothing to start
        return nil
    }

    // Is it a docker service?
    if p.Yaml.Runtime.Image != "" {
        p.docker.name = p.Yaml.Name
        // Source path & mode daemon
        p.docker.opts = []string{"-v", p.Source_path + ":/src/", "-d"}

        // Do we have a socket to prepare?
        if p.Yaml.Runtime.Service.Port == 0 && p.cmd.socket_path != "" {
            if err := p.socket_prepare(); err != nil {
                return err
            }
            p.docker.socket_path = "/forjj-sockets"
            p.docker.opts = append(p.docker.opts, "-v", p.cmd.socket_path+":"+p.docker.socket_path)
        } else {
            return fmt.Errorf("Forjj connect to remote url - Not yet implemented\n")
        }

        // Check if the container exists and is started.
        if _, err := p.docker_container_restart(p.cmd.command, p.cmd.args); err != nil {
            return err
        }

        for i := 1; i < 30; i++ {
            if p.CheckServiceUp() {
                return nil
            }
            time.Sleep(time.Second)
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
    return body == "OK"
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
