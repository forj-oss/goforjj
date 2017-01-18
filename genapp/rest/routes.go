package main

import (
	"github.com/forj-oss/goforjj"
)

var routes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"Quit", "GET", "/quit", Quit},
	Route{"Ping", "GET", "/ping", goforjj.PingHandler},
	Route{"Create", "POST", "/create", Create},
	Route{"Update", "POST", "/update", Update},
	Route{"Maintain", "POST", "/maintain", Maintain},
}
