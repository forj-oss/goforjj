package goforjj

func (f *YamlFlag)found_action_def(action string) (found bool) {
	if f == nil {
		return
	}

	if f.Actions == nil {
		return true
	}
	if found, _ = InArray(action, f.Actions) ; found {
		return
	}
	found, _ = InArray(setup, f.Actions)
	return
}
