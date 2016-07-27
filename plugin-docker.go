package goforjj

import "fmt"

// Function to start an existing container or create and run a new one
func (p *PluginDef) docker_container_restart(cmd string, args []string) (string, error) {
    if p.Yaml.Runtime.Image == "" {
        return "", fmt.Errorf("docker_image is missing in the driver definition. driver ignored.\n")
    }
    ret, err := docker_container_status(p.docker.name)
    if err != nil {
        return "", err
    }
    switch ret {
    case "started":
        return "", nil
    case "":
        dopts := []string{"--name", p.docker.name}
        dopts = append(dopts, p.docker.opts...)
        return docker_container_run(dopts, p.Yaml.Runtime.Image, cmd, args)
    default:
        return docker_container_start(p.docker.name)
    }

}

// Function to run a container
func docker_container_run(docker_opts []string, image, cmd string, args []string) (string, error) {
    cmd_args := append([]string{}, "sudo", "docker", "run")
    cmd_args = append(cmd_args, docker_opts...)
    cmd_args = append(cmd_args, image, cmd)
    cmd_args = append(cmd_args, args...)

    return cmd_run(cmd_args)
}

func docker_container_start(name string) (string, error) {
    cmd_args := append([]string{}, "sudo", "docker", "start", name)
    return cmd_run(cmd_args)
}

func docker_container_status(name string) (string, error) {
    cmd_args := append([]string{}, "sudo", "docker", "inspect", "--format", "{{ .State.Status }}", name)
    return cmd_run(cmd_args)
}

func docker_container_remove(name string) (string, error) {
    cmd_args := append([]string{}, "sudo", "docker", "rm", "-f", name)
    return cmd_run(cmd_args)
}
