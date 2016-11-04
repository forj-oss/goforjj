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
		Source("main.go", 0644, "//", template_shell_main, false).
		Source("app.go", 0644, "//", template_shell_app, true)
}

func (a *App) init_model_rest() {
	a.Models.Create("REST API").
		Source("main.go", 0644, "//", template_rest_main, false).
		Source(a.Yaml.Name+"_plugin.go", 0644, "//", template_rest_plugin, false).
		Source("create.go", 0644, "//", template_rest_create, false).
		Source("update.go", 0644, "//", template_rest_update, false).
		Source("maintain.go", 0644, "//", template_rest_maintain, false).
		Source("app.go", 0644, "//", template_rest_app, false).
		Source("cli.go", 0644, "//", template_rest_cli, false).
		Source("handlers.go", 0644, "//", template_rest_handlers, false).
		Source("actions.go", 0644, "//", template_rest_actions, false).
		Source("routes.go", 0644, "//", template_rest_routes, false).
		Source("router.go", 0644, "//", template_rest_router, false).
		Source("log.go", 0644, "//", template_rest_log, false).
		Source("Dockerfile", 0644, "#", template_rest_dockerfile, false).
		Source("bin/build.sh", 0755, "", template_rest_build, false).
		Source("bin/publish-alltags.sh", 0755, "", template_rest_publish, false).
		Source("ca_certificates/README.md", 0644, "", template_rest_careadme, false).
		Source("README.md", 0644, "", template_rest_readme, false).
		Source("yaml-structs.go", 0644, "//", template_rest_structs, true)
}
