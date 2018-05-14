package goforjj

import (
	"github.com/forj-oss/forjj-modules/trace"
	"os"
)

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
