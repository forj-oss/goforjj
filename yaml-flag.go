package goforjj

// data structure in /objects/<Object Name>/flags/<flag name>
//     flags:
//       <flag name>:
//         help: string - Help attached to the flag
//         required: bool - true if this flag is required.
type YamlFlag struct {
	Options       YamlFlagOptions `yaml:",inline"`
	Help          string
	FormatRegexp  string   `yaml:"format-regexp"`
	Actions       []string `yaml:"only-for-actions"`
	CliCmdActions []string `yaml:"cli-exported-to-actions"`
	Type          string   `yaml:"of-type"`
	FlagScope     string   `yaml:"flag-scope"` // 'object' by default. Flag is not prefixed by instance name.
	// 'instance' Flag is prefixed by instance name if certain condition.
	FieldScope string `yaml:"fields-scope"` // 'object' by default. Means field is added at Object level.
	// 'instance' Means fields is added at object instance level.
	extentSource bool // true if the flag is defined by source as extent. Requires ExtendRelPath from YamlPlugin
}

type YamlFlagOptions struct {
	Required bool
	Hidden   bool   // Used by the plugin.
	Default  string // Used by the plugin.
	Secure   bool   // true if the data must be securely stored, ie not in the git repo. The flag must be defined in 'common' or 'maintain' flag group.
	Envar    string // Environment variable name to use.
}

func (f *YamlFlag) found_action_def(action string) (found bool) {
	if f == nil {
		return
	}

	if f.Actions == nil {
		return true
	}
	if found, _ = InArray(action, f.Actions); found {
		return
	}
	found, _ = InArray(setup, f.Actions)
	return
}

// IsExtentFlag is True if the flag was defined as extent.
func (f *YamlFlag) IsExtentFlag() bool {
	return f.extentSource
}
