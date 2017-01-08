package goforjj

import (
	"encoding/json"
	"fmt"
	"github.com/forj-oss/goforjj/trace"
)

// Function which will execute the action requested.
// If the is a REST API, communicate with REST API protocol
// else start a shell or a container to get the json data.
func (p *PluginDef) PluginRunAction(action string, d PluginReqData) (*PluginResult, error) {
	if p.service {
		return p.api_do(action, d)
	}
	return p.shell_do(action, d)
}

// Internally execute the REST POST Call with parameters
// returns the decoded data into predefined recognized PluginResult sructure
func (p *PluginDef) api_do(action string, d PluginReqData) (*PluginResult, error) {
	p.url.Path = action
	var (
		data []byte
		err  error
	)

	if data, err = json.Marshal(d); err != nil {
		return nil, err
	}

	gotrace.Trace("POST %s with '%s'", p.url.String(), string(data))
	resp, body, errs := p.req.Post(p.url.String()).Send(string(data)).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	gotrace.Trace("Json data returned: \n%s", body)
	var result PluginResult

	if err := json.Unmarshal([]byte(body), &result.Data); err != nil {
		return nil, err
	}

	gotrace.Trace("data extracted: \n%#v", result.Data)

	if result.Data.ErrorMessage != "" {
		result.State_code = resp.StatusCode
		return &result, fmt.Errorf("Plugin issue detected: %s", result.Data.ErrorMessage)
	}
	return &result, nil
}

// Internally execute the command with collection of parameters
// returns the decoded data into predefined recognized PluginResult sructure
func (p *PluginDef) shell_do(action string, d PluginReqData) (*PluginResult, error) {
	cmd_args := []string{action}

	for k, v := range d.Args {
		cmd_args = append(cmd_args, "--"+k)
		cmd_args = append(cmd_args, v)
	}

	if len(d.ReposData) > 0 {
		if dj, err := json.Marshal(d.ReposData); err == nil {
			cmd_args = append(cmd_args, "--data", string(dj))
		}
	}

	body, err := cmd_run(cmd_args)
	if err != nil {
		return nil, err
	}

	var result PluginResult

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
