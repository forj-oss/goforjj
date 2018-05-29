package main

import (
	"net/http"

	"github.com/forj-oss/goforjj"
)

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from '{{.Yaml.Name}}.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int) {
	var p *__MYPLUGIN__Plugin

	// This is where you shoud write your Create code. Following lines are typical code for a basic plugin.
	if pr, ok := req.check_source_existence(ret); !ok {
		return
	} else {
		p = pr
	}

	if !p.initialize_from(req, ret) {
		return
	}

	// Example of the core task (req.a contains the list of args your plugin requires)
	//if ! p.create_jenkins_sources(req.a.ForjjInstanceName, ret) {
	//    return
	//}
	// If your plugin is an upstream plugin, you will need to get the list of Requested Repositories from "req.ReposData"

	if !p.save_yaml(ret, req.Forj.ForjjInstanceName) {
		return
	}
	return
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from '{{.Yaml.Name}}.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {
	var p *__MYPLUGIN__Plugin

	// This is where you shoud write your Update code. Following lines are typical code for a basic plugin.
	if pr, ok := req.check_source_existence(ret); !ok {
		return
	} else {
		p = pr
	}

	instance := req.Forj.ForjjInstanceName
	if !p.load_yaml(ret, instance) {
		return
	}

	if !p.update_from(req, ret) {
		return
	}

	// Example of the core task
	//if ! p.update_jenkins_sources(ret) {
	//    return
	//}

	if !p.save_yaml(ret, instance) {
		return
	}
	return
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from '{{.Yaml.Name}}.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoMaintain(r *http.Request, req *MaintainReq, ret *goforjj.PluginData) (httpCode int) {
	// This is where you shoud write your Maintain code. Following lines are typical code for a basic plugin.
	if !req.check_source_existence(ret) {
		return
	}

	if !req.instantiate(ret) {
		return
	}
	return
}
