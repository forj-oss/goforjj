package goforjj

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/forj-oss/forjj-modules/trace"
	"github.com/parnurzeal/gorequest"
)

const defaultTimeout = 32 * time.Second

func (p *PluginDef) define_as_local_paths() {
	p.SourceMount = p.Source_path
	p.DestMount = p.DeployPath
	if _, err := os.Stat(p.Source_path); err != nil {
		os.MkdirAll(p.Source_path, 0755)
	}

	// Workspace path
	if p.Workspace_path != "" {
		p.WorkspaceMount = p.Workspace_path
	} else {
		p.WorkspaceMount = ""
	}

}

func (p *PluginDef) command_start_service() (err error) {
	if _, err = p.define_socket(); err != nil {
		return
	}

	p.define_as_local_paths()

	cmd_args := []string{p.cmd.command}
	cmd_args = append(cmd_args, p.cmd.args...)
	_, err = cmd_run(cmd_args)
	return
}

// PluginStartService This function start the service as daemon and register it
// If the service is already started, just use it.
func (p *PluginDef) PluginStartService() (err error) {
	if !p.service {
		// Nothing to start
		return nil
	}
	gotrace.Trace("Starting plugin service...")

	switch {
	case p.local_debug: // Local debug Nothing to start
		p.define_as_local_paths()
		gotrace.Trace("Local debugger activated. The service is not started.")
	case p.Yaml.Runtime.Docker.Image != "": // Docker to start
		err = p.docker_start_service()
	default: // Command to start
		err = p.command_start_service()
	}

	if err != nil {
		return
	}

	// Do a ping of the service.
	p.CheckServiceUp()
	return
}

func (p *PluginDef) CheckServiceUp() bool {
	gotrace.Trace("Checking service response.")
	if p.cmd.socket_file != "" {
		if _, err := os.Stat(path.Join(p.cmd.socket_path, p.cmd.socket_file)); os.IsNotExist(err) {
			return false
		}
	}
	gotrace.Trace("Pinging the service.")

	p.define_socket()
	p.url.Path = "ping"
	_, body, err := p.req.Get(p.url.String()).End()
	if err != nil {
		fmt.Printf("Issue to ping the Plugin service '%s'. %s\n", p.Yaml.Name, err)
	}
	p.service_booted = true
	gotrace.Trace("Service is UP.")
	return strings.Trim(body, " \n") == "OK"
}

// Create the socket link for http and his path.
func (p *PluginDef) socket_prepare() (err error) {
	// Define it once
	if p.req != nil {
		return
	}

	p.cmd.socket_file = p.Yaml.Name + ".sock"
	socket := path.Join(p.cmd.socket_path, p.cmd.socket_file)
	p.req = gorequest.New()
	p.req.Transport.Dial = func(_, _ string) (net.Conn, error) {
		return net.DialTimeout("unix", socket, defaultTimeout)
	}
	p.url, err = url.Parse("http://" + p.cmd.socket_file)

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
	if !p.service || !p.service_booted || p.local_debug {
		return
	}
	p.url.Path = "quit"
	p.req.Get(p.url.String()).End()

	if p.Yaml.Runtime.Docker.Image != "" {
		for i := 0; i <= 10; i++ {
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
