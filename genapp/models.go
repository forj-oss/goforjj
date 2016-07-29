package main

import (
    "fmt"
    "os"
)

func (a *App) init_model() {
    switch a.Yaml.Runtime.Service_type {
    case "REST API":
        a.init_model_rest()
    case "shell":
        a.init_model_shell()
    default:
        fmt.Printf("Invalid Service type. Must be 'REST API' or 'shell'.")
        os.Exit(1)
    }
}


// Do not define any plugin.go file, which is the first plugin golang file created by the Forjj plugin creator.
func (a *App) init_model_shell() {
    a.Models.Create("shell").
        Source("main.go", template_shell_main, false).
        Source("app.go", template_shell_app,   true)
}

func (a *App) init_model_rest() {
    a.Models.Create("REST API").
        Source("main.go",         template_rest_main,     false).
        Source("app.go",          template_rest_app,      false).
        Source("cli.go",          template_rest_cli,      false).
        Source("handlers.go",     template_rest_handlers, false).
        Source("actions.go",      template_rest_actions,  false).
        Source("routes.go",       template_rest_routes,   false).
        Source("router.go",       template_rest_router,   false).
        Source("log.go",          template_rest_log,      false).
        Source("yaml-structs.go", template_rest_structs,  true)
}
