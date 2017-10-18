package goforjj

// data structure in /runtime/docker

type DockerStruct struct {
	Image   string // Image name declared by the plugin.
	Dood    bool              `yaml:",omitempty"`
	Volumes []string          `yaml:",omitempty"`
	Env     map[string]string `yaml:",omitempty"`
	User    string            `yaml:",omitempty"`
}
