package goforjj

import (
    "fmt"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
)

// Function to start an existing container or create and run a new one
func (p *PluginDef) docker_container_restart(cmd string, args []string) (string, error) {
    if p.Yaml.Runtime.Image == "" {
        return "", fmt.Errorf("docker_image is missing in the driver definition. driver ignored.\n")
    }
    gotrace.Trace("Restarting container '%s' with action: %s, args: %s", p.docker.name, cmd, args)
    ret, err := docker_container_status(p.docker.name)
    if err != nil {
        return "", err
    }
    switch ret {
    case "started":
        return "", nil
    case "":
        dopts := []string{"--name", p.docker.name}
        p.docker.complete_opts_with(p.docker.volumes, p.docker.envs)
        dopts = append(dopts, p.docker.opts...)
        gotrace.Trace("Booting container '%s' status", p.docker.name)
        return docker_container_run(dopts, p.Yaml.Runtime.Image, cmd, args)
    default:
        gotrace.Trace("Starting container '%s' status", p.docker.name)
        return docker_container_start(p.docker.name)
    }

}

// Function to run a container
func docker_container_run(docker_opts []string, image, cmd string, args []string) (string, error) {
    gotrace.Trace("Starting container from image '%s'", image)
    cmd_args := append([]string{}, "sudo", "docker", "run")
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
    cmd_args := append([]string{}, "sudo", "docker", "stop", name)
    return cmd_run(cmd_args)
}

func docker_container_start(name string) (string, error) {
    gotrace.Trace("Starting container '%s'", name)
    cmd_args := append([]string{}, "sudo", "docker", "start", name)
    return cmd_run(cmd_args)
}

func docker_container_status(name string) (string, error) {
    gotrace.Trace("Checking container '%s' status", name)
    cmd_args := append([]string{}, "sudo", "docker", "inspect", "--format", "{{ .State.Status }}", name)
    return cmd_run(cmd_args)
}

func docker_container_logs(name string) (string, error) {
    gotrace.Trace("Getting container '%s' logs", name)
    cmd_args := append([]string{}, "sudo", "docker", "logs", name)
    return cmd_run(cmd_args)
}

func docker_container_remove(name string) (string, error) {
    gotrace.Trace("Removing container '%s'", name)
    cmd_args := append([]string{}, "sudo", "docker", "rm", "-f", name)
    return cmd_run(cmd_args)
}
