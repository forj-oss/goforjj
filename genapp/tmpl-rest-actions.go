package main

const template_rest_actions = `package main

// You can remove following comments.
// It has been designed fo you, to implement the core of your plugin task.
//
// You can use use it to write your own plugin handler for additional functionnality
// Like Index which currently return a basic code.

import (
    "fmt"
    "os"
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
)

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from '{{.Yaml.Name}}.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(w http.ResponseWriter, r *http.Request, req_data *CreateReq, ret_data *goforjj.PluginData) (httpCode int) {

    // This is where you shoud write your Update code. Following line is for Demo only.
    fmt.Fprintf(os.Stdout,"%#v\n", req_data)

    return
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from '{{.Yaml.Name}}.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(w http.ResponseWriter, r *http.Request, req_data *UpdateReq, ret_data *goforjj.PluginData) (httpCode int){

    // This is where you shoud write your create code. Following line is for Demo only.
    fmt.Fprintf(os.Stdout,"%#v\n", req_data)

    return
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from '{{.Yaml.Name}}.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoMaintain(w http.ResponseWriter, r *http.Request, req_data *MaintainReq, ret_data *goforjj.PluginData) (httpCode int) {

    // This is where you shoud write your Update code. Following line is for Demo only.
    fmt.Fprintf(os.Stdout,"%#v\n", req_data)

    return
}
`
