package goforjj

import (
    "gopkg.in/yaml.v2"
    "fmt"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
    "path"
    "os"
)

// Load yaml raw data in YamlPlugin data structure
func (p *PluginDef) PluginDefLoad(yaml_data []byte) (error) {
    return yaml.Unmarshal([]byte(yaml_data), &p.Yaml)
}

// Initialize Plugin with Definition data.
func (p *PluginDef) PluginInit(instance string) (error) {
    gotrace.Trace("Initializing plugin instance '%s'", instance)
    if p.Yaml.Name == "" {
        return fmt.Errorf("Unable to initialize the plugin without Plugin definition.")
    }
    if err := p.def_runtime_context(); err != nil {
        return err
    }

    // To define a unique container name based on workspace name.
    p.docker.name = instance + "-" + p.Yaml.Name
    gotrace.Trace("Service mode : %t", p.service)
    return nil
}

func (p *PluginDef)def_runtime_context() (error) {
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

// Set plugin source path. Created later by docker_start_service
func (p *PluginDef) PluginSetSource(path string) {
    p.Source_path = path
}

func (p *PluginDef) PluginSetWorkspace(path string) {
    p.Workspace_path = path
}

// Declare the socket path. It will be created later by function socket_prepare
func (p *PluginDef) PluginSocketPath(path string) {
    p.cmd.socket_path = path
}

func (p *PluginDef) PluginDockerBin(thePath string) error {
    if thePath == "" {
        gotrace.Trace("PluginDockerBin : '%s'.", thePath)
        return nil
    }
    if _, err := os.Stat(thePath) ; err == nil {
        p.dockerBin = path.Clean(thePath)
    } else {
        return fmt.Errorf("Invalid PluginDockerBin '%s'. %s", thePath, err)
    }
    return nil
}
