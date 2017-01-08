package goforjj

import (
	"encoding/json"
	"fmt"
	//    "github.com/forj-oss/goforjj/trace"
	//    "os"
	//   "os/exec"
	//    "strings"
	//    "syscall"
)

// Function to print out json data
func (p *PluginResult) JsonPrint() error {
	if b, err := json.Marshal(p); err != nil {
		return err
	} else {
		fmt.Printf("%s\n", b)
	}
	return nil
}

/*
// Load data returned by the plugin in the internal structure of Forjj core.
+func (p *PluginsDefinition)LoadResult(res *PluginResult) error {
+ if p.plugins == nil { p.plugins = make(map[string]PluginDef) }
+ return nil
 }
*/
func (o *PluginData) StatusAdd(n string, args ...interface{}) string {
	if o.Status != "" {
		o.Status += "\n"
	}
	s := fmt.Sprintf(n, args...)
	o.Status += s
	return s
}

func (o *PluginData) Errorf(s string, args ...interface{}) string {
	s = fmt.Sprintf(s, args...)
	o.ErrorMessage = s
	return s
}
