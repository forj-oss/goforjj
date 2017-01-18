package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var cliApp __MYPLUGIN__App

func main() {
	cliApp.init()

	switch kingpin.MustParse(cliApp.App.Parse(os.Args[1:])) {
	case "service start":
		cliApp.start_server()
	default:
		kingpin.Usage()
	}
}
