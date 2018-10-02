package runcontext

import (
	"os"
	"strings"
)

// RunContext define a docker run collection of mount and env which can be re-used for a second container in DooD mode
type RunContext struct {
	volumes          []string
	options          []string
	env              map[string]string
	sharedName       string
	shared           bool
	from             bool
	hasContainerAdds bool
	addVolume        func(string)
	addEnv           func(string, string)
	addHiddenEnv     func(string, string)
	addOptions       func(...string)
}

// NewRunContext create a new docker Run Context to shared
func NewRunContext(sharedName string, volumeSize int) (ret *RunContext) {
	ret = new(RunContext)
	ret.sharedName = sharedName
	if sharedName != "" {
		ret.shared = true
	}
	ret.options = make([]string, 0)
	ret.volumes = make([]string, 0, volumeSize)
	ret.env = make(map[string]string)
	return
}

// DefineContainerFuncs define 3 container functions to update container options automatically.
func (r *RunContext) DefineContainerFuncs(addVolume func(string), addEnv func(string, string), addHiddenEnv func(string, string), addOptions func(...string)) {
	r.addEnv = addEnv
	r.addHiddenEnv = addHiddenEnv
	r.addVolume = addVolume
	r.addOptions = addOptions
	r.hasContainerAdds = true
}

// GetFrom buiild the docker Run Context from the shared name found from GetEnv
func (r *RunContext) GetFrom() (shared bool) {
	if r == nil {
		return
	}
	if r.sharedName != "" {
		v := os.Getenv(r.sharedName)
		if v == "" {
			return
		}
		volume := false
		env := false
		for _, element := range strings.Split(v, " ") {
			if !volume && !env {
				switch element {
				case "-v":
					volume = true
				case "-e":
					env = true
				default:
					r.AddOptions(element)
				}
				continue
			}
			if volume {
				r.AddVolume(element)
				volume = false
				continue
			}

			if element == "" {
				env = false
				continue
			}
			vals := strings.Split(element, "=")
			if len(vals) == 1 {
				r.AddEnv(vals[0], "")
			} else {
				r.AddEnv(vals[0], vals[1])
			}
			env = false
			continue

		}
		shared = true
	}
	r.from = shared
	return
}

// AddVolume add a volume option to the current Context
func (r *RunContext) AddVolume(volume string) *RunContext {
	if r == nil {
		return nil
	}
	r.volumes = append(r.volumes, volume)
	if r.hasContainerAdds {
		r.addVolume(volume)
	}
	return r
}

// AddEnv add a env option to the current Context
func (r *RunContext) AddEnv(key, value string) *RunContext {
	if r == nil {
		return nil
	}
	r.env[key] = value
	if r.hasContainerAdds {
		r.addEnv(key, value)
	}
	return r
}

// AddOptions add a env option to the current Context
func (r *RunContext) AddOptions(options ...string) *RunContext {
	if r == nil {
		return nil
	}
	r.options = append(r.options, options...)
	if r.hasContainerAdds {
		r.addOptions(options...)
	}
	return r
}

// AddFromEnv create a new docker run env (-e) from existing environment variable.
func (r *RunContext) AddFromEnv(key string) *RunContext {
	if r == nil {
		return nil
	}
	if v := strings.Trim(os.Getenv(key), " "); v != "" {
		r.env[key] = ""
		if r.hasContainerAdds {
			r.addHiddenEnv(key, v)
		}
	}
	return r
}

// AddShared return the shared variables with context options if shared == true
func (r *RunContext) AddShared() {
	if r == nil || r.sharedName == "" || r.from || !r.shared || !r.hasContainerAdds {
		return
	}

	sharedValue := r.buildSharedName()

	r.addEnv(r.sharedName, sharedValue)
	return
}

// BuildOptions creates an array of options for cmd.
func (r *RunContext) BuildOptions() (ret []string) {
	if r == nil {
		return
	}

	ret = make([]string, 0, len(r.volumes)*2+len(r.options)+len(r.env)*2)

	for _, volume := range r.volumes {
		ret = append(ret, "-v", volume)
	}

	ret = append(ret, r.options...)

	for key, value := range r.env {
		ret = append(ret, "-e")
		if value == "" {
			ret = append(ret, key)
		} else {
			ret = append(ret, key+"="+value)
		}
	}

	if r.shared {
		ret = append(ret, "-e")
		ret = append(ret, r.sharedName+"="+r.buildSharedName())
	}
	return
}

// buildSharedName creates a shell string to use with `docker run`
func (r *RunContext) buildSharedName() (sharedValue string) {
	values := make([]string, 0, len(r.volumes)+len(r.options)+len(r.env))

	for _, volume := range r.volumes {
		values = append(values, "-v '"+volume+"'")
	}

	for _, option := range r.options {
		values = append(values, "'"+option+"'")
	}

	for key, value := range r.env {
		if value == "" {
			values = append(values, "-e '"+key+"'")
		} else {
			values = append(values, "-e '"+key+"="+value+"'")
		}
	}
	return strings.Join(values, " ")
}
