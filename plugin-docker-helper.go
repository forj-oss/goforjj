package goforjj

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"os"
	"path"
	"strings"
)

// Function to start an existing container or create and run a new one
func (p *PluginDef) docker_container_restart(cmd string, args []string) (string, error) {
	Image := p.Yaml.Runtime.Docker.Image
	if Image == "" {
		return "", fmt.Errorf("runtime/docker/image is missing in the driver definition. driver ignored.\n")
	}
	Image += ":" + p.Version

	// Docker pull policy: Consider latest image tag as Mutable and others as Immutable.
	// Until Docker comes with a docker run --pull ... https://github.com/moby/moby/issues/34394
	// Forjj will do the refresh only for latest image by default.
	if p.Version == Latest { // Check and refresh image if needed.
		gotrace.Trace("Latest image policy check:")
		if _, err := docker_image_pull(Image); err != nil {
			return "", err
		}
		if container_image, err := docker_inspect(p.docker.name, ".Image") ; err == nil && container_image != "" {
			if image_info, err := docker_inspect(container_image, ".RepoTags") ; err != nil {
				return "", err
			} else {
				if ! strings.Contains(image_info, Image) {
					gotrace.Trace("The container '%s' is going to be removed as the image has been updated.",
						p.docker.name)
					if _, err = docker_container_remove(p.docker.name); err != nil {
						return "", err
					}
				} else {
					gotrace.Trace("'%s' do not need to be refreshed.", Image)
				}
			}
		}
	}

	gotrace.Trace("Restarting container '%s' with action: %s, args: %s", p.docker.name, cmd, args)
	ret, err := docker_container_status(p.docker.name)
	if err != nil {
		return "", err
	}
	status := strings.Trim(ret, " \n")
	p.cleanup_socket(status)
	switch status {
	case "running":
		return "", nil
	case "":
		dopts := []string{"--name", p.docker.name}
		p.docker.complete_opts_with(p.docker.volumes, p.docker.envs)
		dopts = append(dopts, p.docker.opts...)
		gotrace.Trace("Booting container '%s' status", p.docker.name)
		return docker_container_run(dopts, Image, cmd, args)
	default:
		gotrace.Trace("Starting container '%s' status", p.docker.name)
		return docker_container_start(p.docker.name)
	}

}

// Function to remove any already existing socket file.
// Usually, needs to be executed if the container is not running.
func (p *PluginDef) cleanup_socket(status string) {
	if status == "running" {
		return
	}
	if p.cmd.socket_file != "" {
		file := path.Join(p.cmd.socket_path, p.cmd.socket_file)
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
			gotrace.Trace("Removed socket file '%s' related to a non running container", file)
		}
	}

}

func docker_container_sudo() (docker_cmd []string) {
	if v := os.Getenv("DOCKER_SUDO") ; v != "" {
		docker_cmd = []string{v}
	} else {
		docker_cmd = []string{}
	}
	return
}

// Function to run a container
func docker_container_run(docker_opts []string, image, cmd string, args []string) (string, error) {
	gotrace.Trace("Starting container from image '%s'", image)
	cmd_args := append(docker_container_sudo(), "docker", "run")
	cmd_args = append(cmd_args, docker_opts...)
	cmd_args = append(cmd_args, image)
	if cmd != "" {
		cmd_args = append(cmd_args, cmd)
	}
	cmd_args = append(cmd_args, args...)

	return cmd_run(cmd_args)
}

func docker_container_stop(name string) (string, error) {
	gotrace.Trace("Stopping container '%s'", name)
	cmd_args := append(docker_container_sudo(), "docker", "stop", name)
	return cmd_run(cmd_args)
}

func docker_container_start(name string) (string, error) {
	gotrace.Trace("Starting container '%s'", name)
	cmd_args := append(docker_container_sudo(), "docker", "start", name)
	return cmd_run(cmd_args)
}

func docker_container_status(name string) (string, error) {
	return docker_inspect(name, ".State.Status")
}

func docker_container_logs(name string) (string, error) {
	gotrace.Trace("Getting container '%s' logs", name)
	cmd_args := append(docker_container_sudo(), "docker", "logs", name)
	return cmd_run(cmd_args)
}

func docker_container_remove(name string) (string, error) {
	gotrace.Trace("Removing container '%s'", name)
	cmd_args := append(docker_container_sudo(), "docker", "rm", "-f", name)
	return cmd_run(cmd_args)
}

func docker_image_pull(name string) (string, error) {
	gotrace.Trace("Pulling image '%s'", name)
	cmd_args := append(docker_container_sudo(), "docker", "pull", name)
	return cmd_run(cmd_args)
}

func docker_inspect(name, data string) (string, error) {
	gotrace.Trace("Getting info '%s' from '%s'", data, name)
	cmd_args := append(docker_container_sudo(), "docker", "inspect", "--format", "{{ " + data + " }}", name)
	return cmd_run(cmd_args)
}
