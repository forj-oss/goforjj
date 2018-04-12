package goforjj

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
