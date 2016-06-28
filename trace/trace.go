package gotrace

import (
        "fmt"
        "runtime"
)

var debug bool

func SetDebug() {
 debug = true
}

func Trace(s string, a ...interface{}) {
 if ! debug { return }

 pc := make([]uintptr, 10)  // at least 1 entry needed
 runtime.Callers(2, pc)
 f := runtime.FuncForPC(pc[0])
 fmt.Printf(fmt.Sprintf("DEBUG %s: %s", f.Name(), s), a ...)
}


