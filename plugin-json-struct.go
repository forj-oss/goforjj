package goforjj

//***************************************
// JSON data structure for shell type of plugin.

type PluginRepo struct {
    Name     string   // name of the repository
    Upstream string   // upstream url
}

type PluginService struct {
    Urls map[string]string
}

// REST API json data
type PluginData struct {
    Repos         map[string]PluginRepo // List of repository data
    Services      PluginService         // web service url. ex: https://github.hpe.com
    Status        string                // Status message
    CommitMessage string                // Action commit message for Create/Update
    ErrorMessage  string                // Found only if error detected
    Files         []string              // List of files managed by the plugin
}

// Shell json data
type PluginResult struct {
    Data          PluginData
    State_code    int   // 200 OK
}
