package goforjj

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/forj-oss/forjj-modules/trace"
)

type DockerService struct {
	Volumes map[string]byte
	Env     map[string]byte
}

func (p *PluginDef) define_socket() (remote bool, err error) {
	if p.Yaml.Runtime.Service.Port == 0 && p.cmd.socket_path != "" {
		err = p.socket_prepare()
		return
	}

	err = fmt.Errorf("Forjj connect to remote url - Not yet implemented\n")
	remote = true
	return
}

// docker_start_service Define how to start
func (p *PluginDef) docker_start_service() (err error) {
	gotrace.Trace("Starting it as docker container '%s'", p.docker.name)

	// initialize
	p.docker.init()

	// mode daemon
	p.docker.opts = append(p.docker.opts, "-d")

	// Source path
	if _, err := os.Stat(p.Source_path); err != nil {
		os.MkdirAll(p.Source_path, 0755)
	}
	p.SourceMount = "/src/"
	p.docker.add_volume(p.Source_path + ":" + p.SourceMount)

	if p.DeployPath != "" { // For compatibility reason with old forjj.
		p.DestMount = "/deploy/"
		p.docker.add_volume(p.DeployPath + ":" + p.DestMount)
	}

	// Workspace path
	if p.Workspace_path != "" {
		p.WorkspaceMount = "/workspace/"
		p.docker.add_volume(p.Workspace_path + ":" + p.WorkspaceMount)
	}

	// Define the socket
	remote_url := false
	remote_url, err = p.define_socket()
	if err != nil {
		return
	}
	if !remote_url {
		p.docker.socket_path = "/tmp/forjj-socks"
		p.docker.add_volume(p.cmd.socket_path + ":" + p.docker.socket_path)
	}

	if p.Yaml.Runtime.Docker.Volumes != nil {
		for _, v := range p.Yaml.Runtime.Docker.Volumes {
			p.docker.add_volume(v)
		}
	}

	if p.Yaml.Runtime.Docker.Env != nil {
		for k, v := range p.Yaml.Runtime.Docker.Env {
			if env := os.ExpandEnv(v); v != env && env != "" {
				gotrace.Trace("expand and set %s from %s to %s", k, v, env)
				p.docker.add_env(k, env)
			} else {
				gotrace.Trace("set %s to %s", k, v)
				p.docker.add_env(k, v)
			}
		}
	}

	if p.Yaml.Runtime.Docker.Dood {
		if p.dockerBin == "" {
			err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path", p.docker.name)
			return
		}
		gotrace.Trace("Adding docker dood information...")
		// TODO: download bin version of docker and mount it, or even communicate with the API directly in the plugin container (go: https://github.com/docker/engine-api)

		if dood_mt_opts, dood_bc_opts, err := p.GetDockerDoodParameters(); err != nil {
			return err
		} else {
			p.docker.opts = append(p.docker.opts, dood_mt_opts...)
			p.docker.opts = append(p.docker.opts, dood_bc_opts...)
		}
	} else {
		p.docker.opts = append(p.docker.opts, "-u", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}

	// Check if the container exists and is started.
	// TODO: Be able to interpret some variables written in the <plugin>.yaml and interpreted here to start the daemon correctly.
	// Ex: all p.cmd_data .* in a golang template would give {{ .socket_path }}, etc...
	if _, err = p.docker_container_restart(p.Yaml.Runtime.Service.Command, p.Yaml.Runtime.Service.Parameters); err != nil {
		return
	}

	err = p.check_service_ready()

	return
}

// GetDockerDoodParameters returns 2 Arrays of strings
//
// mount : contains mount information to share.
// become: contains information required by a container to add user and become that user to start the process
//         Depending on the container, `become` can be ignored if the container do not need to become user
//         sent by forjj.
func (p *PluginDef) GetDockerDoodParameters() (mount, become []string, err error) {
	if !p.Yaml.Runtime.Docker.Dood {
		return nil, nil, fmt.Errorf("Dood not defined by the plugin. Required to use it.")
	}
	// In context of dood, the container must respect few things:
	// - The container is started as root
	// - the container start/entrypoint must grab the UID/GID environment sent by forjj to set the appropriate
	//   unprivileged user. ie useradd or equivalent.
	// - The plugin MUST be executed with UID/GID user context. You can use either su, sudo, or any other user account
	//   substitute.
	//   ie su - <User>
	// - Usually the container should have access to a /bin/docker binary compatible with host docker version.
	//   provided by forjj with --docker-exe
	// - forjj will mount /var/run/docker.sock to /var/run/docker.sock root access limited, no shared group. so you
	//   must use a sudoers so your plugin user could call docker against the host server socket.
	if p.dockerBin == "" {
		err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path", p.docker.name)
		return
	}

	if v := strings.Trim(os.Getenv("DOCKER_DOOD"), " "); v != "" {
		mount = strings.Split(v, " ")
	} else {
		mount = make([]string, 0, 8)
		mount = append(mount, "-v", "/var/run/docker.sock:/var/run/docker.sock")
		mount = append(mount, "-v", p.dockerBin+":/bin/docker")
		mount = append(mount, "-e", "DOOD_SRC="+p.Source_path)
		mount = append(mount, "-e", "DOOD_DEPLOY="+p.DeployPath)
	}

	if v := strings.Trim(os.Getenv("DOCKER_DOOD_BECOME"), ""); v != "" {
		become = strings.Split(v, " ")
	} else {
		become = make([]string, 0, 6)
		become = append(become, "-u", "root:root")
		become = append(become, "-e", "UID="+strconv.Itoa(os.Getuid()))
		become = append(become, "-e", "GID="+strconv.Itoa(os.Getgid()))
	}

	return
}

// Regularly testing the service response. fails after a timeout.
func (p *PluginDef) check_service_ready() (err error) {
	gotrace.Trace("Checking service status...")
	for i := 1; i < 30; i++ {
		time.Sleep(time.Second)

		out := ""
		if out, err = docker_container_status(p.docker.name); err != nil {
			return
		}

		if strings.Trim(out, " \n") != "running" {
			out, err = docker_container_logs(p.docker.name)
			if err == nil {
				out = fmt.Sprintf("docker logs:\n---\n%s---\n", out)
			} else {
				out = fmt.Sprintf("%s\n", err)
			}
			docker_container_remove(p.docker.name)
			err = fmt.Errorf("%sContainer '%s' has stopped unexpectedely.", out, p.Yaml.Name)
			return
		} else {
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
	if d.envs == nil {
		d.init()
	}
	if ok, _ := regexp.MatchString("^.*(:.*(:(ro|rw))?)?$", volume); ok {
		d.volumes[volume] = 'v'
	}
}

func (d *docker_container) add_env(key, value string) {
	if d.envs == nil {
		d.init()
	}
	env := key + "=" + value
	d.envs[env] = 'e'
}

func (d *docker_container) complete_opts_with(val ...map[string]byte) {
	// Calculate the expected array size
	tot := len(d.opts)

	for _, v := range val {
		tot += len(v) * 2
	}

	// allocate
	r := make([]string, 0, tot)

	// append
	r = append(r, d.opts...)
	for _, v := range val {
		for k, o := range v {
			r = append(r, "-"+string(o))
			r = append(r, k)
		}
	}
	d.opts = r
}
