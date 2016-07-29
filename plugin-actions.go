package goforjj

import (
    "encoding/json"
    "github.hpe.com/christophe-larsonneur/goforjj/trace"
)

// Function which will execute the action requested.
// If the is a REST API, communicate with REST API protocol
// else start a shell or a container to get the json data.
func (p *PluginDef) PluginRunAction(action string, args map[string]string) (*PluginResult, error) {
    if p.service {
        return p.api_do(action, args)
    }
    return p.shell_do(action, args)
}

// Internally execute the REST POST Call with parameters
// returns the decoded data into predefined recognized PluginResult sructure
func (p *PluginDef) api_do(action string, args map[string]string) (*PluginResult, error) {
    p.url.Path = action
    var (
        data []byte
        err error
    )

    if data, err = json.Marshal(args); err != nil {
        return nil, err
    }

    gotrace.Trace("POST %s with '%s'", p.url.String(), string(data))
    _, body, errs := p.req.Post(p.url.String()).Send(string(data)).End()
    if len(errs) > 0 {
        return nil, errs[0]
    }

    var result PluginResult

    if err := json.Unmarshal([]byte(body), &result); err != nil {
        return nil, err
    }

    return &result, nil
}

// Internally execute the command with collection of parameters
// returns the decoded data into predefined recognized PluginResult sructure
func (p *PluginDef) shell_do(action string, args map[string]string) (*PluginResult, error) {
    cmd_args := []string{action}

    for k, v := range args {
        cmd_args = append(cmd_args, "--" + k)
        cmd_args = append(cmd_args, v)
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
