package goforjj

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/forj-oss/forjj-modules/trace"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const defaultTimeout = 32 * time.Second

// This function start the service as daemon and register it
// If the service is already started, just use it.
func (p *PluginDef) PluginStartService(instance_name string) error {
	if !p.service {
		// Nothing to start
		return nil
	}
	gotrace.Trace("Starting plugin service...")
	// Is it a docker service?
	if is_docker, err := p.docker_start_service(instance_name); is_docker {
		return err
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
	gotrace.Trace("Checking service response.")
	if p.cmd.socket_file != "" {
		if _, err := os.Stat(path.Join(p.cmd.socket_path, p.cmd.socket_file)); os.IsNotExist(err) {
			return false
		}
	}
	gotrace.Trace("Pinging the service.")
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
	if !p.service || !p.service_booted {
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
