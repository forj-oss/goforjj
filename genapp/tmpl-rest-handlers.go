package main

const template_rest_handlers = `package main

// This file has been created by "go generate" as initial code. go generate will never update it, except if you remove it.

// So, update it for your need.

import (
    "fmt"
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "encoding/json"
)

// Index
func Index(w http.ResponseWriter, _ *http.Request) {
    fmt.Fprintf(w ,"FORJJ - {{.Yaml.Name}} driver for FORJJ. It is Implemented as a REST API.")
}

// Do creating route
func Create(w http.ResponseWriter, r *http.Request) {
    var data goforjj.PluginData
    // Create the github.yaml source file.
    // See goforjj/plugin-json-struct.go for json data structure recognized by forjj.

    if err := json.NewEncoder(w).Encode(data); err != nil {
        panic(err)
    }
}

// Do updating route
func Update(w http.ResponseWriter, r *http.Request) {
    var data goforjj.PluginData
    // Update the github.yaml source file.
    // See goforjj/plugin-json-struct.go for json data structure recognized by forjj.

    if err := json.NewEncoder(w).Encode(data); err != nil {
        panic(err)
    }
}

// Do maintaining route
func Maintain(w http.ResponseWriter, r *http.Request) {

}
`
