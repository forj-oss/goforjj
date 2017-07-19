package goforjj

func NewRepo() *PluginRepo {
	r := new(PluginRepo)
	r.Remotes = make(map[string]PluginRepoRemoteUrl)
	r.BranchConnect = make(map[string]string)
	return r
}

// GetUpstream Currently get the 'upstream' if exist or 'origin' url
func (r *PluginRepo) GetUpstream(forGit bool) string {
	if r == nil {
		return ""
	}
	if s, found := r.Remotes["upstream"] ; found {
		if forGit {
			return s.Ssh
		} else {
			return s.Url
		}
	}
	if s, found := r.Remotes["origin"] ; found {
		if forGit {
			return s.Ssh
		} else {
			return s.Url
		}
	}
	return ""
}

func (r *PluginRepo) GetOrigin(forGit bool) string {
	if r == nil {
		return ""
	}
	if s, found := r.Remotes["origin"] ; found {
		if forGit {
			return s.Ssh
		} else {
			return s.Url
		}
	}
	return ""
}

