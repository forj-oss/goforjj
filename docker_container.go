package goforjj

import (
	"fmt"
	"regexp"

	"github.com/forj-oss/forjj-modules/trace"
)

// DockerContainer is the named container information
type DockerContainer struct {
	name        string
	opts        []string
	image       string
	socket_path string
	volumes     map[string]byte
	envs        map[string]byte
	dockerCmd   commandRun
	outFunc     func(line string)
	errFunc     func(line string)
}

// Init initialize the DockerContainer structure.
func (d *DockerContainer) Init() {
	d.opts = make([]string, 0, 5)
	d.volumes = make(map[string]byte)
	d.envs = make(map[string]byte)
	d.dockerCmd.Init(dockerCmd()...)
	d.outFunc = func(line string) {
		fmt.Printf(line)
	}
	d.errFunc = func(line string) {
		fmt.Printf(line)
	}
}

// AddOpts add some docker options (passed before the image name)
func (d *DockerContainer) AddOpts(opts ...string) {
	if d == nil {
		return
	}

	d.opts = append(d.opts, opts...)
}

// SetName defines the docker container name
func (d *DockerContainer) SetName(name string) {
	if d == nil {
		return
	}

	d.name = name
}

// Name return the container name to use.
func (d *DockerContainer) Name() string {
	if d == nil {
		return ""
	}

	return d.name
}

// AddVolume add a volume for a docker container
func (d *DockerContainer) AddVolume(volume string) {
	if d == nil {
		return
	}

	if d.envs == nil {
		d.Init()
	}
	if ok, _ := regexp.MatchString("^.*(:.*(:(ro|rw))?)?$", volume); ok {
		d.volumes[volume] = 'v'
	}
}

// AddEnv add an environment variable to the container via -e key=value
func (d *DockerContainer) AddEnv(key, value string) {
	if d == nil {
		return
	}
	if d.envs == nil {
		d.Init()
	}
	env := key + "=" + value
	d.envs[env] = 'e'
}

// AddHiddenEnv add an environment variable to the container via -e key and the environment variable to the docker run command.
func (d *DockerContainer) AddHiddenEnv(key, value string) {
	if d == nil {
		return
	}
	if d.envs == nil {
		d.Init()
	}
	env := key
	d.envs[env] = 'e'
	d.dockerCmd.AddEnv(key, value)
}

// SetImageName define the image name to use for the container.
func (d *DockerContainer) SetImageName(image string) {
	if d == nil {
		return
	}
	d.image = image
}

func (d *DockerContainer) complete_opts_with(val ...map[string]byte) {
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
func (d *DockerContainer) Run(dockerOpts []string, cmd string, args []string) error {
	gotrace.Trace("Starting container from image '%s'", d.image)

	d.dockerCmd.SetArgs(d.configureDockerCli("run", dockerOpts, cmd, args))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// configureDockerCli Helps to configure docker cli call
func (d *DockerContainer) configureDockerCli(action string, dopts []string, command string, cmdArgs []string) (args []string) {
	args = make([]string, 0, 3+len(dopts)+len(d.opts)+len(cmdArgs))
	args = append(args, action)
	if d.opts != nil {
		args = append(args, d.opts...)
	}
	if dopts != nil {
		args = append(args, dopts...)
	}
	args = append(args, d.image, command)
	if cmdArgs != nil {
		args = append(args, cmdArgs...)
	}
	return
}

// Stop stop the named container
func (d *DockerContainer) Stop() error {
	gotrace.Trace("Stopping container '%s'", d.name)
	d.dockerCmd.SetArgs(d.configureDockerCli("stop", nil, d.Name(), nil))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Start the named container
func (d *DockerContainer) Start() error {
	gotrace.Trace("Starting container '%s'", d.name)
	d.dockerCmd.SetArgs(d.configureDockerCli("start", nil, d.Name(), nil))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Status get container status field
func (d *DockerContainer) Status() (string, error) {
	return d.Inspect(d.name, ".State.Status")
}

// Logs printout the log to the out function.
func (d *DockerContainer) Logs(out func(string)) error {
	gotrace.Trace("Getting container '%s' logs", d.Name())
	d.dockerCmd.SetArgs(d.configureDockerCli("logs", nil, d.Name(), nil))
	return d.dockerCmd.runFlow(out, d.errFunc)
}

// Remove a container
func (d *DockerContainer) Remove() error {
	gotrace.Trace("Removing container '%s'", d.name)
	d.dockerCmd.SetArgs(d.configureDockerCli("rm", []string{"-f"}, d.name, nil))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Pull the container image
func (d *DockerContainer) Pull() error {
	gotrace.Trace("Pulling image '%s'", d.image)
	d.dockerCmd.SetArgs(d.configureDockerCli("pull", nil, d.image, nil))
	return d.dockerCmd.runFlow(d.outFunc, d.errFunc)
}

// Inspect the named container.
func (d *DockerContainer) Inspect(name string, data string) (ret string, _ error) {
	gotrace.Trace("Getting info '%s' from '%s'", data, name)
	d.dockerCmd.SetArgs(d.configureDockerCli("inspect", []string{"--format", "{{ " + data + " }}"}, name, nil))
	return ret, d.dockerCmd.runFlow(func(line string) {
		if ret == "" {
			ret = line
		} else {
			ret += "\n" + line
		}
	}, d.errFunc)
}

