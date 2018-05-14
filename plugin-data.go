package goforjj

import (
	"encoding/json"
	"fmt"
)

// REST API json data
type PluginData struct {
	Repos         map[string]PluginRepo     `json:",omitempty"` // List of repository data
	Services      PluginService             `json:",omitempty"` // web service url. ex: https://github.hpe.com
	Status        string                    // Status message
	CommitMessage string                    `json:",omitempty"` // Action commit message for Create/Update
	ErrorMessage  string                    // Found only if error detected
	Files         map[string][]string       `json:",omitempty"` // List of files managed by the plugin
	Options       map[string]PluginOption   `json:",omitempty"` // List of options needed at maintain use case, returned from create/update. Usually used to provide credentials.
}

const (
	FilesSource = "source"
	FilesDeploy = "deploy"
)

// AddFile add a file the list of files Forjj will take care in GIT.
//
func (d *PluginData) AddFile(where, file string) {
	if where == "" {
		where = FilesDeploy
	}
	if where != FilesDeploy && where != FilesSource {
		where = FilesDeploy
	}
	if d.Files == nil {
		d.Files = make(map[string][]string)
	}
	if v, found := d.Files[where]; !found {
		v = make([]string, 1, 5)
		v[0] = file
		d.Files[where] = v
	} else {
		v = append(v, file)
		d.Files[where] = v
	}
}


// JsonPrint to print out json data
func (p *PluginResult) JsonPrint() error {
	if b, err := json.Marshal(p); err != nil {
		return err
	} else {
		fmt.Printf("%s\n", b)
	}
	return nil
}

// StatusAdd Add status information to the API caller.
func (o *PluginData) StatusAdd(n string, args ...interface{}) string {
	if o.Status != "" {
		o.Status += "\n"
	}
	s := fmt.Sprintf(n, args...)
	o.Status += s
	return s
}

// Errorf to store error message made by all other functions and return it to the 
// API caller
func (o *PluginData) Errorf(s string, args ...interface{}) string {
	s = fmt.Sprintf(s, args...)
	o.ErrorMessage = s
	return s
}
