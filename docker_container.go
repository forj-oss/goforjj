package goforjj

import (
	"fmt"
	"regexp"

	"github.com/forj-oss/forjj-modules/trace"
)

/*
This module is a docker helper for most of forjj command.
It is not a generic docker wrapper. If you want to call docker, I suggest to use the docker library directly

In the future, I expect to remove this code and replace by the docker library directly.
*/

// DockerContainer is the named container information
type DockerContainer struct {
	name        string
	opts        []string // docker RUN Options
	image       string
	socket_path string
	volumes     map[string]byte
	envs        map[string]byte
	dockerCmd   commandRun
	outFunc     func(line string)
	errFunc     func(line string)
}

// Init initialize the DockerContainer structure.
func (d *DockerContainer) Init(name string) {
	if d == nil {
		return
	}

	d.opts = make([]string, 2)
	d.opts[0] = "--name"
	d.opts[1] = name
	d.name = name
	d.volumes = make(map[string]byte)
	d.envs = make(map[string]byte)
	d.dockerCmd.Init(dockerCmd()...)
	d.outFunc = func(line string) {
		fmt.Print(line, "\n")
	}
	d.errFunc = func(line string) {
		fmt.Print(line, "\n")
	}
}

// AddOpts add some docker options (passed before the image name)
func (d *DockerContainer) AddOpts(opts ...string) {
	if d == nil || d.envs == nil {
		return
	}

	d.opts = append(d.opts, opts...)
}

// Name return the container name to use.
func (d *DockerContainer) Name() string {
	if d == nil || d.envs == nil {
		return ""
	}

	return d.name
}

// AddVolume add a volume for a docker container
func (d *DockerContainer) AddVolume(volume string) {
	if d == nil || d.envs == nil {
		return
	}
	if ok, _ := regexp.MatchString("^.*(:.*(:(ro|rw))?)?$", volume); ok {
		d.volumes[volume] = 'v'
	}
}

// AddEnv add an environment variable to the container via -e key=value
func (d *DockerContainer) AddEnv(key, value string) {
	if d == nil || d.envs == nil {
		return
	}
	env := key + "=" + value
	d.envs[env] = 'e'
}

// AddHiddenEnv add an environment variable to the container via -e key and the environment variable to the docker run command.
func (d *DockerContainer) AddHiddenEnv(key, value string) {
	if d == nil || d.envs == nil {
		return
	}
	env := key
	d.envs[env] = 'e'
	d.dockerCmd.AddEnv(key, value)
}

// SetImageName define the image name to use for the container.
func (d *DockerContainer) SetImageName(image string) {
	if d == nil || d.envs == nil {
		return
	}
	d.image = image
}

// complete_opts_with add volumes and environments to the docker run opts
func (d *DockerContainer) complete_opts_with() {
	val := []map[string]byte{d.volumes, d.envs}
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

// Run the container
func (d *DockerContainer) Run(cmd string, args []string) error {
	gotrace.Trace("Starting container from image '%s'", d.image)

	d.dockerCmd.SetArgs(d.configureDockerRunCli("run", d.opts, d.image, cmd, args))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// configureDockerCli Helps to configure docker cli run call
//
// - action then
// - action options then
// - object then
// - command then
// - commands arg
//
// Following docker actions follow this format.
// docker run : https://docs.docker.com/engine/reference/commandline/run/
func (d *DockerContainer) configureDockerRunCli(action string, dopts []string, image, command string, cmdArgs []string) (args []string) {
	args = make([]string, 0, 3+len(dopts)+len(cmdArgs))
	args = append(args, action)
	if dopts != nil {
		args = append(args, dopts...)
	}
	args = append(args, d.image)
	if command != "" {
		args = append(args, command)
	}
	if cmdArgs != nil {
		args = append(args, cmdArgs...)
	}
	return
}

// configureDockerObjectCli
//
// - action then
// - action options then
// - object
//
// Following docker actions follow this format
// docker logs : https://docs.docker.com/engine/reference/commandline/logs/
// docker pull : https://docs.docker.com/engine/reference/commandline/pull/
func (d *DockerContainer) configureDockerObjectCli(action string, dopts []string, object string) []string {
	return d.configureDockerObjectsCli(action, dopts, object)
}

// configureDockerObjectsCli Helps to configure docker cli run call
//
// - action then
// - action options then
// - objects
//
// Following docker actions follow this format
// docker stop    : https://docs.docker.com/engine/reference/commandline/stop/
// docker start   : https://docs.docker.com/engine/reference/commandline/start/
// docker inspect : https://docs.docker.com/engine/reference/commandline/inspect/
func (d *DockerContainer) configureDockerObjectsCli(action string, dopts []string, objects ...string) (args []string) {
	args = make([]string, 0, 1+len(dopts)+len(objects))
	args = append(args, action)
	if dopts != nil {
		args = append(args, dopts...)
	}
	args = append(args, objects...)
	return
}

// Stop stop the named container
func (d *DockerContainer) Stop(dopts []string) error {
	gotrace.Trace("Stopping container '%s'", d.name)
	d.dockerCmd.SetArgs(d.configureDockerObjectsCli("stop", dopts, d.Name()))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Start the named container
func (d *DockerContainer) Start(dopts []string) error {
	gotrace.Trace("Starting container '%s'", d.name)
	d.dockerCmd.SetArgs(d.configureDockerObjectsCli("start", dopts, d.Name()))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Status get container status field
func (d *DockerContainer) Status() (string, error) {
	return d.Inspect(d.name, ".State.Status")
}

// Logs printout the log to the out function.
func (d *DockerContainer) Logs(dopts []string, out func(string)) error {
	gotrace.Trace("Getting container '%s' logs", d.Name())
	d.dockerCmd.SetArgs(d.configureDockerObjectCli("logs", dopts, d.Name()))
	return d.dockerCmd.runFlow(out, d.errFunc)
}

// Remove a container
func (d *DockerContainer) Remove() error {
	gotrace.Trace("Removing container '%s'", d.name)
	d.dockerCmd.SetArgs(d.configureDockerObjectCli("rm", []string{"-f"}, d.name))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Pull the container image
func (d *DockerContainer) Pull(dopts []string) error {
	gotrace.Trace("Pulling image '%s'", d.image)
	d.dockerCmd.SetArgs(d.configureDockerObjectCli("pull", dopts, d.image))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Inspect the named container.
func (d *DockerContainer) Inspect(name string, data string) (ret string, _ error) {
	gotrace.Trace("Getting info '%s' from '%s'", data, name)
	d.dockerCmd.SetArgs(d.configureDockerObjectsCli("inspect", []string{"--format", "{{ " + data + " }}"}, name))
	return ret, d.dockerCmd.runFlow(func(line string) {
		if ret == "" {
			ret = line
		} else {
			ret += "\n" + line
		}
	}, d.errFunc)
}
