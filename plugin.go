package goforjj

import (
    "encoding/json"
    "fmt"
//    "github.hpe.com/christophe-larsonneur/goforjj/trace"
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
