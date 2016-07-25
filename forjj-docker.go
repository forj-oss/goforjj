package goforjj

import (
    "fmt"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
    "os/exec"
    "strings"
)

// Function to start a container in Daemon
func docker_container_run(name, image, action string, docker_opts []string, opts []string) ([]byte, error) {
    if image == "" {
        return []byte{}, fmt.Errorf("docker_image is missing in the driver definition. driver ignored.\n")
    }
    cmd_args := append([]string{}, "sudo", "docker", "run", "-d", "--name", name)
    cmd_args = append(cmd_args, docker_opts...)
    cmd_args = append(cmd_args, image, action)
    cmd_args = append(cmd_args, opts...)

    gotrace.Trace("STARTING: %s\n\n", strings.Join(cmd_args, " "))
    if ret, err := exec.Command(cmd_args[0], cmd_args[1:]...).Output() ; err != nil {
        return []byte{}, err
    } else {
        return ret, nil
    }
}

func docker_container_stat(name string) ([]byte, error) {
    cmd_args := append([]string{}, "sudo", "docker", "ps", "-a", "-f", "{{ .State.Status }}", name)
    gotrace.Trace("STARTING: %s\n\n", strings.Join(cmd_args, " "))
    if ret, err := exec.Command(cmd_args[0], cmd_args[1:]...).Output() ; err != nil {
        return []byte{}, err
    } else {
        return ret, nil
    }
}

func docker_container_remove(name string) ([]byte, error) {
    cmd_args := append([]string{}, "sudo", "docker", "rm", "-f", name)
    gotrace.Trace("STARTING: %s\n\n", strings.Join(cmd_args, " "))
    if ret, err := exec.Command(cmd_args[0], cmd_args[1:]...).Output() ; err != nil {
        return []byte{}, err
    } else {
        return ret, nil
    }
}
