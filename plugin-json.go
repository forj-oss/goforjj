package goforjj

func (d *PluginData)AddFile(file string){
    if d.Files == nil {
        d.Files = make([]string, 0, 5)
    }
    d.Files = append(d.Files, file)
}
