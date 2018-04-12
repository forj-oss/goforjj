package goforjj

const (
	source = "source"
	deploy = "deploy"
)

// AddFile add a file the list of files Forjj will take care in GIT.
//
func (d *PluginData) AddFile(where, file string) {
	if where == "" {
		where = deploy
	}
	if where != deploy && where != source {
		where = deploy
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
