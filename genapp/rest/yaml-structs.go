package main

type CreateReq struct {
	Forj    map[string]string
	Objects CreateArgReq
}

type CreateArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $Opts := .Yaml.Objects }}\
	repo map[string]map[string]string // __MYPLUGIN:     {{ go_vars $Objectname}} map[string]map[string]string `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{ end }}\
}

type UpdateReq struct {
	Forj    map[string]string
	Objects UpdateArgReq
}

type UpdateArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $Opts := .Yaml.Objects }}\
	repo map[string]map[string]string // __MYPLUGIN:     {{ go_vars $Objectname}} map[string]map[string]string `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{ end }}\
}

type MaintainReq struct {
	Forj    map[string]string
	Objects MaintainArgReq
}

type MaintainArgReq struct {
	// __MYPLUGIN: {{ range $Objectname, $Opts := .Yaml.Objects }}\
	// __MYPLUGIN:     {{ go_vars $Objectname}} map[string]map[string]string `json:"{{$Objectname}}"` // Object details
	// __MYPLUGIN: {{ end }}\
}
