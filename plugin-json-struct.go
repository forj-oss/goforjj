package goforjj

//***************************************
// JSON data structure for shell type of plugin.

type PluginRepo struct {
    Name     string   // name of the repository
    Upstream string   // upstream url TODO: obsolete. To remove shorlty in plugins.
    Exist    bool     // True is the repo exist.
    Remotes map[string]string        // k: remote name, v: remote url
    BranchConnect map[string]string  // k: local branch name, v: remote/branch
    Flow string       // Information given by forjj when it ask to apply some fixes on the repo (upstream changes)
}

type PluginService struct {
    Urls map[string]string
}

// REST API json data
type PluginData struct {
    Repos         map[string]PluginRepo   // List of repository data
    Services      PluginService           // web service url. ex: https://github.hpe.com
    Status        string                  // Status message
    CommitMessage string                  // Action commit message for Create/Update
    ErrorMessage  string                  // Found only if error detected
    Files         []string                // List of files managed by the plugin
    Options       map[string]PluginOption // List of options needed at maintain use case. Usually used to provide credentials.
}

type PluginOption struct {
    Help string   // Help about plugin options required at maintain phase
    Value string  // Value set/loaded at create/update phase
}

// Shell json data
type PluginResult struct {
    Data          PluginData
    State_code    int   // 200 OK
}
