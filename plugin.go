package goforjj

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/forj-oss/forjj-modules/trace"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/yaml.v2"
)

const Latest = "latest"

type Driver struct {
	// Driver define an instance of a driver
	Result         *PluginResult         // Json data structured returned.
	Yaml           *YamlPlugin           // Yaml data definition
	name           string                // container name
	Source_path    string                // Plugin source path from Forjj point of view
	Workspace_path string                // Plugin Workspace path from Forjj point of view
	DeployPath     string                // Plugin Deployment path
	DeployName     string                // Plugin Deployment name in path
	service        bool                  // True if the service is started as daemon
	service_booted bool                  // True if the service is started
	container      DockerContainer       // Define data to start the plugin as docker container
	cmd            commandRun            // Define data to start the service process
	req            *gorequest.SuperAgent // REST API request
	url            *url.URL              // REST API url
	dockerBin      string                // Docker Path Binary to a docker binary to mount in a dood container.
	SourceMount    string                // Where the driver will have his source code.
	DestMount      string                // Where the driver will have his generated code.
	WorkspaceMount string                // where the driver has his workspace.
	Version        string                // Plugin version to load
	key64          string                // Base64 symetric key
	local_debug    bool                  // true to bypass starting container or binary. Expect it be started in a running
	// instance of the driver from a debugger
	sourceDefPath string // Path to the source file to complete driver definition
	// Loaded in
}

const defaultTimeout = 32 * time.Second

func NewDriver(plugin *YamlPlugin) (p *Driver) {
	p = new(Driver)
	p.Yaml = plugin
	return
}

// PluginDefLoad Load yaml raw data in YamlPlugin data structure
func (p *Driver) PluginDefLoad(yaml_data []byte) error {

	return yaml.Unmarshal([]byte(yaml_data), p.Yaml)
}

// PluginInit Initialize Plugin with Definition data.
func (p *Driver) PluginInit(instance string) error {
	gotrace.Trace("Initializing plugin instance '%s'", instance)
	if p.Yaml.Name == "" {
		return fmt.Errorf("Unable to initialize the plugin without Plugin definition.")
	}
	if err := p.def_runtime_context(); err != nil {
		return err
	}

	// To define a unique container name based on workspace name.
	p.name = instance + "-" + p.Yaml.Name
	gotrace.Trace("Service mode : %t", p.service)
	return nil
}

func (p *Driver) RunningFromDebugger() {
	p.local_debug = true
}

// PluginSetSource Set plugin source path. Created later by docker_start_service
func (p *Driver) PluginSetSource(path string) {
	p.Source_path = path
}

func (p *Driver) PluginSetWorkspace(path string) {
	p.Workspace_path = path
}

func (p *Driver) PluginSetDeployment(path string) {
	p.DeployPath = path
}

func (p *Driver) PluginSetDeploymentName(name string) {
	p.DeployName = name
}

// PluginSocketPath Declare the socket path. It will be created later by function socket_prepare
func (p *Driver) PluginSocketPath(path string) {
	p.cmd.socket_path = path
}

func (p *Driver) PluginDockerBin(thePath string) error {
	if thePath == "" {
		gotrace.Trace("PluginDockerBin : '%s'.", thePath)
		return nil
	}
	// Check in case of paths like "/something/~/something/"
	if thePath[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Unable to get Current USER information. %s", err)
		}
		dir := usr.HomeDir
		thePath = filepath.Join(dir, thePath[2:])
	}
	if _, err := os.Stat(thePath); err == nil {
		p.dockerBin = path.Clean(thePath)
	} else {
		return fmt.Errorf("Invalid PluginDockerBin '%s'. %s", thePath, err)
	}
	return nil
}

func (p *Driver) PluginSetVersion(version string) {
	if version == "" {
		p.Version = Latest
	} else {
		p.Version = version
	}
	gotrace.Trace("Plugin version selected: %s", p.Version)
}

// PluginLoadFrom do a load of the plugin Def Runtime section
// This information is saved by forjj to avoid reloding the plugin.yaml
// A plugin already loaded is not refreshed.
// NOTE: Workspace_path, Source_path and SourceMount must be set in PluginDef to make it work.
// TODO: Add a Plugin refresh? Not sure if forjj could do it or not differently...
func (p *Driver) PluginLoadFrom(name string, runtime *YamlPluginRuntime) error {
	if name == "" || runtime == nil {
		return fmt.Errorf("Internal Error: PluginRuntimeReloadFrom: name cannot be empty and plugin cannot be nil")
	}
	if p.Yaml.Name != "" {
		gotrace.Trace("'%s' is not loaded from the workspace cache.", p.Yaml.Name)
		return nil
	}
	p.Yaml.Name = name

	p.Yaml.Runtime = *runtime
	gotrace.Trace("'%s' has been reloaded.", p.Yaml.Name)
	return nil
}

// PluginRunAction Function which will execute the action requested.
// If the plugin is a REST API, communicate with real basic REST API protocol
// Basic RESTFul means : GET/POST, simple unique route, no version, payload with everything.
// If needed in a next iteration, we can move the API to match fully the RESTFul API with forjj objects/actions.
// else start a shell or a container to get the json data.
func (p *Driver) PluginRunAction(action string, d *PluginReqData) (*PluginResult, error) {
	p.url.Path = action
	var (
		data []byte
		err  error
	)

	if data, err = json.Marshal(d); err != nil {
		return nil, err
	}

	jsonData, _ := json.MarshalIndent(d, "", "  ")

	p.define_socket()

	gotrace.Trace("POST %s with '%s'", p.url.String(), string(jsonData))
	resp, body, errs := p.req.Post(p.url.String()).Send(string(data)).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var result PluginResult

	if err := json.Unmarshal([]byte(body), &result.Data); err != nil {
		return nil, err
	}

	if dataDisplayed, err := json.MarshalIndent(&result.Data, "", "  "); err == nil {
		gotrace.Trace("data returned: \n%s", string(dataDisplayed))
	} else {
		gotrace.Trace("data returned: \n%#v", result.Data)
	}

	if result.Data.ErrorMessage != "" {
		result.State_code = resp.StatusCode
		return &result, nil
	}
	return &result, nil
}

// GetDockerDoodParameters returns 2 Arrays of strings
//
// mount : contains mount information to share.
// become: contains information required by a container to add user and become that user to start the process
//         Depending on the container, `become` can be ignored if the container do not need to become user
//         sent by forjj.
func (p *Driver) GetDockerDoodParameters() (mount, become []string, err error) {
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
		err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path", p.container.Name())
		return
	}

	if v := strings.Trim(os.Getenv("DOCKER_DOOD"), " "); v != "" {
		mount = strings.Split(v, " ")
	} else {
		mount = make([]string, 0, 8)
		mount = append(mount, "-v", "/var/run/docker.sock:/var/run/docker.sock")
		mount = append(mount, "-v", p.dockerBin+":/bin/docker")
		mount = append(mount, "-e", "DOOD_SRC="+p.Source_path)
		mount = append(mount, "-e", "DOOD_DEPLOY="+path.Join(p.DeployPath, p.DeployName))
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

// PluginStartService This function start the service as daemon and register it
// If the service is already started, just use it.
func (p *Driver) PluginStartService() (err error) {
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

func (p *Driver) CheckServiceUp() bool {
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

// PluginStopService To stop the plugin service if the service was started before by goforjj
func (p *Driver) PluginStopService() {
	if p == nil || !p.service || !p.service_booted || p.local_debug {
		return
	}
	p.url.Path = "quit"
	p.req.Get(p.url.String()).End()

	if p.Yaml.Runtime.Docker.Image != "" {
		for i := 0; i <= 10; i++ {
			time.Sleep(time.Second)
			if out, _ := p.container.Status(); out != "started" {
				return
			}
		}
		if out, _ := p.container.Status(); out == "started" {
			p.container.Stop(nil)
		}
	}
}

// ServiceAddEnv add environment variable to the service runner
func (p *Driver) ServiceAddEnv(name, value string, hidden bool) {

}

// --------------- Internal functions

func (p *Driver) def_runtime_context() error {
	switch p.Yaml.Runtime.Service_type {
	case "REST API": // REST API Service started as daemon
		p.service = true

	case "shell": // Shell/json process
		p.service = false
	default:
		return fmt.Errorf("Error! Invalid '%s' service_type. Supports only 'REST API' and 'shell'. Use shell as default.", p.Yaml.Runtime.Service_type)
	}
	return nil
}

// Function to start an existing container or create and run a new one
func (p *Driver) docker_container_restart(cmd string, args []string) error {
	Image := p.Yaml.Runtime.Docker.Image
	if Image == "" {
		return fmt.Errorf("runtime/docker/image is missing in the driver definition. driver ignored")
	}
	Image += ":" + p.Version

	// Docker pull policy: Consider latest image tag as Mutable and others as Immutable.
	// Until Docker comes with a docker run --pull ... https://github.com/moby/moby/issues/34394
	// Forjj will do the refresh only for latest image by default.
	if p.Version == Latest { // Check and refresh image if needed.
		gotrace.Trace("Latest image policy check:")
		if err := p.container.Pull(nil); err != nil {
			return err
		}
		if container_image, err := p.container.Inspect(p.container.Name(), ".Image"); err == nil && container_image != "" {
			if image_info, err := p.container.Inspect(container_image, ".RepoTags"); err != nil {
				return err
			} else {
				if !strings.Contains(image_info, Image) {
					gotrace.Trace("The container '%s' is going to be removed as the image has been updated.",
						p.container.Name())
					if err = p.container.Remove(); err != nil {
						return err
					}
				} else {
					gotrace.Trace("'%s' do not need to be refreshed.", Image)
				}
			}
		}
	}

	gotrace.Trace("Restarting container '%s' with action: %s, args: %s", p.container.Name(), cmd, args)
	ret, _ := p.container.Status()
	status := strings.Trim(ret, " \n")
	p.cleanup_socket(status)
	switch status {
	case "running":
		return nil
	case "":
		gotrace.Trace("Booting container '%s' status", p.container.Name())
		return p.container.Run(cmd, args)
	default:
		gotrace.Trace("Starting container '%s' status", p.container.Name())
		return p.container.Start(nil)
	}

}

// Function to remove any already existing socket file.
// Usually, needs to be executed if the container is not running.
func (p *Driver) cleanup_socket(status string) {
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

func (p *Driver) define_socket() (remote bool, err error) {
	if p.Yaml.Runtime.Service.Port == 0 && p.cmd.socket_path != "" {
		err = p.socket_prepare()
		return
	}

	err = fmt.Errorf("Forjj connect to remote url - Not yet implemented\n")
	remote = true
	return
}

// docker_start_service Define how to start
func (p *Driver) docker_start_service() (err error) {
	gotrace.Trace("Starting it as docker container '%s'", p.container.Name())

	// Initialize forjj plugins docker container.
	p.container.Init(p.DeployName + "-" + p.name)

	Image := p.Yaml.Runtime.Docker.Image
	if Image == "" {
		return fmt.Errorf("runtime/docker/image is missing in the driver definition. driver ignored")
	}
	Image += ":" + p.Version
	p.container.SetImageName(Image)

	// mode daemon
	p.container.AddOpts("-d")

	// Source path
	if _, err := os.Stat(p.Source_path); err != nil {
		os.MkdirAll(p.Source_path, 0755)
	}
	p.SourceMount = "/src/"
	p.container.AddVolume(p.Source_path + ":" + p.SourceMount)

	if p.DeployPath != "" { // For compatibility reason with old forjj.
		p.DestMount = "/deploy/"
		p.container.AddVolume(p.DeployPath + ":" + p.DestMount)
	}

	// Workspace path
	if p.Workspace_path != "" {
		p.WorkspaceMount = "/workspace/"
		p.container.AddVolume(p.Workspace_path + ":" + p.WorkspaceMount)
	}

	// Define the socket
	remote_url := false
	remote_url, err = p.define_socket()
	if err != nil {
		return
	}
	if !remote_url {
		p.container.socket_path = "/tmp/forjj-socks"
		p.container.AddVolume(p.cmd.socket_path + ":" + p.container.socket_path)
	}

	if p.Yaml.Runtime.Docker.Volumes != nil {
		for _, v := range p.Yaml.Runtime.Docker.Volumes {
			p.container.AddVolume(v)
		}
	}

	if p.Yaml.Runtime.Docker.Env != nil {
		for k, v := range p.Yaml.Runtime.Docker.Env {
			if env := os.ExpandEnv(v); v != env && env != "" {
				gotrace.Trace("expand and set %s from %s to %s", k, v, env)
				p.container.AddEnv(k, env)
			} else {
				gotrace.Trace("set %s to %s", k, v)
				p.container.AddEnv(k, v)
			}
		}
	}

	if p.key64 != "" {
		p.container.AddHiddenEnv("FORJJ_KEY", p.key64)
	}

	if p.Yaml.Runtime.Docker.Dood {
		if p.dockerBin == "" {
			err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path", p.container.Name())
			return
		}
		gotrace.Trace("Adding docker dood information...")
		// TODO: download bin version of docker and mount it, or even communicate with the API directly in the plugin container (go: https://github.com/docker/engine-api)

		if dood_mt_opts, dood_bc_opts, err := p.GetDockerDoodParameters(); err != nil {
			return err
		} else {
			p.container.AddOpts(dood_mt_opts...)
			p.container.AddOpts(dood_bc_opts...)
		}
	} else {
		p.container.AddOpts("-u", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}

	p.container.complete_opts_with()

	// Check if the container exists and is started.
	// TODO: Be able to interpret some variables written in the <plugin>.yaml and interpreted here to start the daemon correctly.
	// Ex: all p.cmd_data .* in a golang template would give {{ .socket_path }}, etc...
	if err = p.docker_container_restart(p.Yaml.Runtime.Service.Command, p.Yaml.Runtime.Service.Parameters); err != nil {
		return
	}

	err = p.check_service_ready()

	return
}

// check_service_ready Regularly testing the service response. fails after a timeout.
func (p *Driver) check_service_ready() (err error) {
	gotrace.Trace("Checking service status...")
	for i := 1; i < 30; i++ {
		time.Sleep(time.Second)

		out := ""
		if out, err = p.container.Status(); err != nil {
			return
		}

		if strings.Trim(out, " \n") != "running" {
			err = p.container.Logs(nil, func(line string) {
				if out == "" {
					out = line
				} else {
					out += "\n" + line
				}
			})
			if err == nil {
				out = fmt.Sprintf("docker logs:\n---\n%s---\n", out)
			} else {
				out = fmt.Sprintf("%s\n", err)
			}
			p.container.Remove()
			err = fmt.Errorf("%sContainer '%s' has stopped unexpectedely.", out, p.Yaml.Name)
			return
		} else {
			return
		}

	}
	err = fmt.Errorf("Plugin Service '%s' not started successfully as docker container '%s'. check docker logs\n", p.Yaml.Name, p.container.Name())
	return
}

func (p *Driver) define_as_local_paths() {
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

func (p *Driver) command_start_service() (err error) {
	if _, err = p.define_socket(); err != nil {
		return
	}

	p.define_as_local_paths()

	cmd_args := p.cmd.command
	cmd_args = append(cmd_args, p.cmd.args...)
	_, err = cmd_run(cmd_args)
	return
}

// Create the socket link for http and his path.
func (p *Driver) socket_prepare() (err error) {
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
