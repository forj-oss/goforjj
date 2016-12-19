package main

import (
	"github.com/forjj-oss/goforjj"
	"log"
	"os"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) check_source_existence(ret *goforjj.PluginData) (status bool) {
	log.Print("Checking __MYPLUGINNAME__ source code path existence.")

	if _, err := os.Stat(r.Forj["forjj-source-mount"]); err == nil {
		ret.Errorf("Unable to maintain __MYPLUGINNAME__ instances. '%s' is inexistent or innacessible.\n",
			r.Forj["forjj-source-mount"])
		return
	}
	ret.StatusAdd("environment checked.")
	return true
}

func (r *MaintainReq) instantiate(ret *goforjj.PluginData) (status bool) {

	return true
}
