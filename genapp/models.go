package main

import (
	"fmt"
	"os"
)

const rest = "rest"

func (a *App) init_model() (model string) {
	switch a.Yaml.Runtime.Service_type {
	case "REST API":
		return a.init_model_rest()
	default:
		fmt.Print("Invalid Service type. Must be 'REST API' or 'shell'.")
		os.Exit(1)
	}
	return
}

func (a *App) init_model_rest() string {
	a.Models.Create(rest, a.template_path).
		Source("main.go", 0644, "//", "main.go", false).
		Source(a.Yaml.Name+"_plugin.go", 0644, "//", "plugin.go", false).
		Source("create.go", 0644, "//", "create.go", false).
		Source("update.go", 0644, "//", "update.go", false).
		Source("maintain.go", 0644, "//", "maintain.go", false).
		Source("app.go", 0644, "//", "app.go", false).
		Source("cli.go", 0644, "//", "cli.go", false).
		Source("handlers.go", 0644, "//", "handlers.go", false).
		Source("actions.go", 0644, "//", "actions.go", false).
		Source("routes.go", 0644, "//", "routes.go", false).
		Source("router.go", 0644, "//", "router.go", false).
		Source("log.go", 0644, "//", "log.go", false).
		Source("Dockerfile", 0644, "#", "Dockerfile", false).
		Source("bin/build.sh", 0755, "", "build.sh", false).
		Source("bin/publish-alltags.sh", 0755, "", "publish.sh", false).
		Source("ca_certificates/README.md", 0644, "", "careadme.md", false).
		Source("README.md", 0644, "", "rest-readme.md", false).
		Source("yaml-structs.go", 0644, "//", "yaml-structs.go", true)
	return rest
}
