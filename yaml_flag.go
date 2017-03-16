package goforjj

func (f *YamlFlag)found_action_def(action string) (found bool) {
	if f == nil {
		return
	}

	if f.Actions == nil {
		return true
	}
	if found, _ = in_array(action, f.Actions) ; found {
		return
	}
	found, _ = in_array(setup, f.Actions)
	return
}
