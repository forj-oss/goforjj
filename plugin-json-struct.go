package goforjj

//***************************************
// JSON data structure of plugin output.

type PluginRepo struct {
	Name          string            // name of the repository
	Exist         bool              // True if the repo exist.
	Remotes       map[string]string // k: remote name, v: remote url
	BranchConnect map[string]string // k: local branch name, v: remote/branch
}

type PluginService struct {
	Urls map[string]string
}

type PluginRepoData struct {
	Templates []string          // RepoTemplates to apply
	Title     string            // Repo Description
	Users     map[string]string // Users and rights given
	Groups    map[string]string // Groups and rights given
	Flow      string            // Flow applied to the Repo
	Instance  string            // Instance managing the upstream repo.
	Options   map[string]string `yaml:",omitempty"` // More options forjj can send out as defined by the plugin itself.
}

// REST API json data
type PluginData struct {
	Repos         map[string]PluginRepo     `json:",omitempty"` // List of repository data
	ReposData     map[string]PluginRepoData `json:",omitempty"` // Data associated to each Repository that the plugin should manage.
	Services      PluginService             `json:",omitempty"` // web service url. ex: https://github.hpe.com
	Status        string                    // Status message
	CommitMessage string                    `json:",omitempty"` // Action commit message for Create/Update
	ErrorMessage  string                    // Found only if error detected
	Files         []string                  `json:",omitempty"` // List of files managed by the plugin
	Options       map[string]PluginOption   `json:",omitempty"` // List of options needed at maintain use case, returned from create/update. Usually used to provide credentials.
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

//***************************************
// JSON data structure of plugin input.
// See plugin-actions.go about how those structs are managed.

type PluginReqData struct {
	Args      map[string]string         `json:"args"`       // Collection of <plugin>.yaml arguments (create/update/maintain) to communicate to the plugin.
	ReposData map[string]PluginRepoData `json:",omitempty"` // Data associated to each Repository that the plugin should manage.
}
