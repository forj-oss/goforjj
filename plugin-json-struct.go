package goforjj

//***************************************
// JSON data structure of plugin output.

type PluginRepo struct {
	Name          string                         // name of the repository
	Exist         bool                           // True if the repo exist.
	Remotes       map[string]PluginRepoRemoteUrl // k: remote name, v: remote url
	BranchConnect map[string]string              // k: local branch name, v: remote/branch
	Owner         string `json:",omitempty"`     // Owner name return by the plugin.
}

type PluginRepoRemoteUrl struct {
	Url string // Public URL (http or https)
	Ssh string // SSH String formatted as (ssh://User@Server:Path or User@Server:Path) for GIT
}

type PluginService struct {
	Urls map[string]string
}



type PluginOption struct {
	Help  string // Help about plugin options required at maintain phase
	Value string // Value set/loaded at create/update phase
}

// Shell json data
type PluginResult struct {
	Data       PluginData
	State_code int // 200 OK
}
