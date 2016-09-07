package goforjj

func (r *PluginRepoData)Initialize_Options(a ...string) {
    for _, d := range a {
        r.Options[d] = ""
    }
}

func (r *PluginRepoData)SetDefaults(defaults map[string]string) {
    if v, found := defaults["flow"] ; found && v!= "" && r.Flow == "" {
        r.Flow = v
    }
    if v, found := defaults["instance"] ; found && v!= "" && r.Instance == "" {
        r.Instance = v
    }
    for k,d := range r.Options {
        if v, found := defaults[k] ; found && v!= "" && d == "" {
            r.Options[k] = v
        }
    }
}

func (r *PluginRepoData)UpdateFrom(source *PluginRepoData) {
    if source.Title != "" {
        r.Title = source.Title
    }
    if source.Users != nil && len(source.Users) >0 {
        r.Users = source.Users
    }
    if source.Flow != "" {
        r.Flow = source.Flow
    }
    if source.Groups != nil && len(source.Groups) >0 {
        r.Groups = source.Groups
    }
    if source.Instance != "" {
        r.Instance = source.Instance
    }
    for k,d := range source.Options {
        if _, found := r.Options[k] ; found {
            r.Options[k] = d
        }
    }
}
