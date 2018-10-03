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
	"syscall"
	"time"

	"github.com/forj-oss/goforjj/runcontext"

	"github.com/forj-oss/forjj-modules/trace"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/yaml.v2"
)

const SocketPathLimit = 108 // See syscall.RawSockaddrUnix.Path in ztypes_linux_amd64.go - Linux limit

const Latest = "latest"

type Driver struct {
	// Driver define an instance of a driver
	Result         *PluginResult         // Json data structured returned.
	Yaml           *YamlPlugin           // Yaml data definition
	name           string                // container name
	service        bool                  // True if the service is started as daemon
	service_booted bool                  // True if the service is started
	container      DockerContainer       // Define data to start the plugin as docker container
	cmd            commandRun            // Define data to start the service process
	req            *gorequest.SuperAgent // REST API request
	url            *url.URL              // REST API url
	dockerBin      string                // Docker Path Binary to a docker binary to mount in a dood container.

	// in docker run syntax, -v BasePath:BaseMount
	basePath  string
	baseMount string

	// in docker run syntax, -v Source_path:SourceMount
	Source_path string // Plugin source path from Forjj point of view
	SourceMount string // Where the driver will have his source code.

	// in docker run syntax, -v Workspace_path:WorkspaceMount
	Workspace_path string // Plugin Workspace path from Forjj point of view
	WorkspaceMount string // where the driver has his workspace.

	// in docker run syntax, -v DeployPath:DestMount
	DeployPath string // Plugin Deployment path
	DestMount  string // Where the driver will have his generated code.

	DeployName string // Plugin Deployment name in path

	Version     string // Plugin version to load
	key64       string // Base64 symetric key
	local_debug bool   // true to bypass starting container or binary. Expect it be started in a running
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

// PluginSetSource Set plugin source path from forjj perspective. Created later by docker_start_service
func (p *Driver) PluginSetSource(path string) {
	p.Source_path = path
}

// PluginSetSourceMount Set plugin source path mount where source will be mounted in the plugin container.
func (p *Driver) PluginSetSourceMount(path string) {
	p.SourceMount = path
}

// PluginSetWorkspaceMount set workspace path from forjj perspective
func (p *Driver) PluginSetWorkspaceMount(path string) {
	p.WorkspaceMount = path
}

// PluginSetDeploymentMount set Deploy path from forjj perspective
func (p *Driver) PluginSetDeploymentMount(path string) {
	p.DestMount = path
}

// PluginSetWorkspace set workspace path from forjj perspective
func (p *Driver) PluginSetWorkspace(path string) {
	p.Workspace_path = path
}

// PluginSetDeployment set Deploy path from forjj perspective
func (p *Driver) PluginSetDeployment(path string) {
	p.DeployPath = path
}

// PluginBase define the source Base mount to use for DooD mount
func (p *Driver) PluginBase(mount string) {
	paths := strings.Split(mount, ":")
	if len(paths) < 2 {
		return
	}
	p.basePath = path.Clean(paths[0])
	p.baseMount = path.Clean(paths[1])
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

	if os.Getenv("DOCKER_DOOD") != "" {
		gotrace.Info("DooD context: workspace 'docker-bin-path' setup is ignored.")
		return nil
	}

	// Check in case of paths like '^~/'
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

func (p *Driver) checkDockerDooD() (err error) {
	if !p.Yaml.Runtime.Docker.Dood {
		return fmt.Errorf("Dood not defined by the plugin. Required to use it")
	}
	return
}

// getDockerDooDGroupID determine the Docker Group ID of /var/run/docker.sock
func (p *Driver) getDockerDooDGroupID() (dockerGrpID uint32, err error) {
	if v := strings.Trim(os.Getenv("DOCKER_DOOD_GROUP"), " "); v != "" {
		if i, convErr := strconv.ParseInt(v, 10, 32); err != nil {
			err = fmt.Errorf("Unable to convert '%s' defined by %s. %s", v, "DOCKER_DOOD_GROUP", convErr)
		} else {
			dockerGrpID = uint32(i)
		}
	} else {
		if s, err := os.Stat("/var/run/docker.sock"); err != nil {
			return 0, err
		} else {
			if v, ok := s.Sys().(*syscall.Stat_t); ok {
				dockerGrpID = v.Gid
			}
		}
	}
	return
}

// Deprecated: GetDockerDoodParameters is kept for build compatibility. It is replaced by DefineDockerDooD()
func (p *Driver) GetDockerDoodParameters() (mount, become []string, err error) {
	if err = p.checkDockerDooD(); err != nil {
		return
	}
	// In context of dood, the container must respect few things:
	// - The container is started as root
	// - the container start/entrypoint must grab the UID/GID environment sent by forjj to set the appropriate
	//   unprivileged user. ie useradd or equivalent.
	// - The plugin MUST be executed with UID/GID user context. So, the plugin container entrypoint should use either su, sudo, or any other user account
	//   substitute to become and start the plugin process.
	//   ie su - <User>
	// - Usually the container should have access to a /bin/docker binary compatible with host docker version.
	//   provided by forjj with --docker-exe-path or workspace docker-bin-path
	// - forjj will mount /var/run/docker.sock to /var/run/docker.sock root access limited, with shared group.
	//   To run the docker against this socket, your entrypoint must have a docker group with same id as docker.sock host.

	// TODO: Ignore this step if docker have to use a tcp connection instead.

	var dockerGrpID uint32

	dockerGrpID, err = p.getDockerDooDGroupID()
	if err != nil {
		return
	}

	dockerDooD := runcontext.NewRunContext("DOCKER_DOOD", 12)
	dockerDooD, err = p.defineDockerDooD(dockerDooD, dockerGrpID)
	if err != nil {
		return
	}
	mount = dockerDooD.BuildOptions()

	become = p.defineDockerDooDBecome(runcontext.NewRunContext("DOCKER_DOOD_BECOME", 6)).BuildOptions()
	return
}

// defineDockerDooD detect and build the DOCKER_DOOD setup
func (p *Driver) defineDockerDooD(dockerDooD *runcontext.RunContext, dockerGrpID uint32) (ret *runcontext.RunContext, err error) {
	ret = dockerDooD
	if !dockerDooD.GetFrom() {
		if p.dockerBin == "" {
			err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path or setup in 'forjj workspace docker-bin-path', or DOCKER_DOOD is empty", p.container.Name())
			return
		}

		dockerDooD.AddVolume("/var/run/docker.sock:/var/run/docker.sock").
			AddVolume(p.dockerBin+":/bin/docker").
			AddEnv("DOOD_SRC", p.Source_path).
			AddEnv("DOOD_DEPLOY", path.Join(p.DeployPath, p.DeployName))
		if dockerGrpID != 0 {
			dockerDooD.AddEnv("DOCKER_DOOD_GROUP", fmt.Sprintf("%d", dockerGrpID))
		}
	}
	return
}

// defineDockerDooDBecome detect and build the DOCKER_DOOD_BECOME setup
func (p *Driver) defineDockerDooDBecome(dockerDooDBecome *runcontext.RunContext) *runcontext.RunContext {
	if !dockerDooDBecome.GetFrom() {
		dockerDooDBecome.AddOptions("-u", "root:root").
			AddEnv("UID", strconv.Itoa(os.Getuid())).
			AddEnv("GID", strconv.Itoa(os.Getgid()))
	}
	return dockerDooDBecome
}

// DefineDockerDood detect and/or define DooD required parameters if the plugin requires it.
//
// It uses goforjj/runcontext module to define and share for new DooD containers the DooD setup.
//
// It manages 2 different context:
// - DOCKER_DOOD. It regroups options to enable DooD with docker socket, docker static binary and docker group ID
//   Those data are set in the new container that forjj will create thanks to addVolume/Env/Opts functions given
//   and "DOCKER_DOOD" will be the last docker env variable which contains all docker run options for the same,
//   shared in the new container to be started by forjj.
//   It requires the container to start as root , in order to update/create the docker group in the container if missing
//   if the container have to be used as a user (non root) it must be created/assigned in the docker image.
//
// - DOCKER_DOOD_BECOME. It regroups options to enable impersonation in the container. The container started as root
//   will ask the container to update few things and become the wanted user with a specific UID/GID given.
//   The container have to update the wanted user UID and GID
func (p *Driver) DefineDockerDood() (err error) {
	if err = p.checkDockerDooD(); err != nil {
		return
	}
	// In context of dood, the container must respect few things:
	// - The container is started as root
	// - the container start/entrypoint must grab the UID/GID environment sent by forjj to set the appropriate
	//   unprivileged user. ie useradd or equivalent.
	// - The plugin MUST be executed with UID/GID user context. So, the plugin container entrypoint should use either su, sudo, or any other user account
	//   substitute to become and start the plugin process.
	//   ie su - <User>
	// - Usually the container should have access to a /bin/docker binary compatible with host docker version.
	//   provided by forjj with --docker-exe-path or workspace docker-bin-path
	// - forjj will mount /var/run/docker.sock to /var/run/docker.sock root access limited, with shared group.
	//   To run the docker against this socket, your entrypoint must have a docker group with same id as docker.sock host.

	// TODO: Ignore this step if docker have to use a tcp connection instead.

	var dockerGrpID uint32

	dockerGrpID, err = p.getDockerDooDGroupID()
	if err != nil {
		return
	}

	dockerDooD := runcontext.NewRunContext("DOCKER_DOOD", 12)
	dockerDooD.DefineContainerFuncs(p.container.AddVolume, p.container.AddEnv, p.container.AddHiddenEnv, p.container.AddOpts)
	dockerDooD, err = p.defineDockerDooD(dockerDooD, dockerGrpID)
	if err != nil {
		return
	}
	dockerDooD.AddShared()

	dockerDooDBecome := runcontext.NewRunContext("DOCKER_DOOD_BECOME", 6)
	dockerDooDBecome.DefineContainerFuncs(p.container.AddVolume, p.container.AddEnv, p.container.AddHiddenEnv, p.container.AddOpts)
	p.defineDockerDooDBecome(dockerDooDBecome)
	dockerDooDBecome.AddShared()

	return
}

// DefineDockerProxyParameters return the list of Proxy parameters
// Shared as DOCKER_DOOD_PROXY
func (p *Driver) DefineDockerProxyParameters() {
	if p == nil {
		return
	}

	dockerDooDProxy := runcontext.NewRunContext("DOCKER_DOOD_PROXY", 6)
	dockerDooDProxy.DefineContainerFuncs(p.container.AddVolume, p.container.AddEnv, p.container.AddHiddenEnv, p.container.AddOpts)
	if !dockerDooDProxy.GetFrom() {
		dockerDooDProxy.AddFromEnv("https_proxy").
			AddFromEnv("http_proxy").
			AddFromEnv("no_proxy")
	}
	dockerDooDProxy.AddShared()
	return
}

// SetDefaultMounts defines container (src/deploy/workspace) mounts to default path
func (p *Driver) SetDefaultMounts() {
	if p == nil {
		return
	}
	p.DestMount = "/deploy/"
	p.SourceMount = "/src/"
	p.WorkspaceMount = "/workspace/"
}

// DefineDockerForjjMounts create a share of forjj driver mounts
func (p *Driver) DefineDockerForjjMounts() error {
	if p == nil {
		return fmt.Errorf("driver is nil")
	}
	if p.SourceMount == "" || p.DestMount == "" || p.WorkspaceMount == "" {
		return fmt.Errorf("Container mounts not set")
	}
	srcContext := runcontext.NewRunContext("DOOD_SOURCE", 12)
	srcContext.DefineContainerFuncs(p.container.AddVolume, p.container.AddEnv, p.container.AddHiddenEnv, p.container.AddOpts)

	// Source path
	if _, err := os.Stat(p.Source_path); err != nil {
		os.MkdirAll(p.Source_path, 0755)
	}
	srcContext.AddVolume(p.Source_path + ":" + p.SourceMount)

	if p.DeployPath != "" { // For compatibility reason with old forjj.
		srcContext.AddVolume(p.DeployPath + ":" + p.DestMount)
	}

	// Workspace path
	if p.Workspace_path != "" {
		srcContext.AddVolume(p.Workspace_path + ":" + p.WorkspaceMount)
	}
	srcContext.AddShared()
	return nil
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
	if p.container.ContainerHasChanged() {
		p.container.Remove()
		status = ""
	}
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
	// Initialize forjj plugins docker container.
	p.container.Init(p.DeployName + "-" + p.name)

	gotrace.Trace("Starting it as docker container '%s'", p.container.Name())

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

	if err = p.DefineDockerForjjMounts() ; err != nil {
		return
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

	p.DefineDockerProxyParameters()

	if p.Yaml.Runtime.Docker.Dood {
		if p.dockerBin == "" {
			err = fmt.Errorf("Unable to activate Dood on docker container '%s'. Missing --docker-exe-path", p.container.Name())
			return
		}
		gotrace.Trace("Adding docker dood information...")
		// TODO: download bin version of docker and mount it, or even communicate with the API directly in the plugin container (go: https://github.com/docker/engine-api)

		if err := p.DefineDockerDood(); err != nil {
			return err
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
	if s := len(socket); s >= SocketPathLimit {
		// Eliminate the Invalid Argument standard message due to this linux socket limit.
		return fmt.Errorf("Socket path exceed linux array size limit (%d). socket '%s' length is %d", SocketPathLimit, socket, s)
	}
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
