package goforjj

import (
	"regexp"
)

type DockerService struct {
	Volumes map[string]byte
	Env     map[string]byte
}

func (d *docker_container) init() {
	d.opts = make([]string, 0, 5)
	d.volumes = make(map[string]byte)
	d.envs = make(map[string]byte)
}

func (d *docker_container) add_volume(volume string) {
	if d.envs == nil {
		d.init()
	}
	if ok, _ := regexp.MatchString("^.*(:.*(:(ro|rw))?)?$", volume); ok {
		d.volumes[volume] = 'v'
	}
}

func (d *docker_container) add_env(key, value string) {
	if d.envs == nil {
		d.init()
	}
	env := key + "=" + value
	d.envs[env] = 'e'
}

func (d *docker_container) complete_opts_with(val ...map[string]byte) {
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
