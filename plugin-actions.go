package goforjj

import (
	"encoding/json"
	"fmt"
	"github.com/forj-oss/goforjj/trace"
)

// PluginRunAction Function which will execute the action requested.
// If the plugin is a REST API, communicate with real basic REST API protocol
// Basic RESTFul means : GET/POST, simple unique route, no version, payload with everything.
// If needed in a next iteration, we can move the API to match fully the RESTFul API with forjj objects/actions.
// else start a shell or a container to get the json data.
func (p *PluginDef) PluginRunAction(action string, d PluginReqData) (*PluginResult, error) {
	return p.api_do(action, d)
}

// api_do Internally execute the REST POST Call with parameters
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
