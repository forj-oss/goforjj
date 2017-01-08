package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func (a *__MYPLUGIN__App) start_server() {

	a.server_set()
	Start := true // Move server to up status

	server_chan := make(chan bool, 1)
	for {
		int_sig := make(chan os.Signal, 1)
		signal.Notify(int_sig, syscall.SIGINT, syscall.SIGTERM)

		ln, err := net.Listen("unix", a.socket)
		if err != nil {
			log.Fatal("listen error:", err)
		}

		// Interruption handler.
		go func() {
			log.Print("Interruption handler started.")
			if _, ok := <-int_sig; !ok {
				log.Print("Exiting interruption handler...")
				signal.Stop(int_sig)
				return // Interruption handler aborted.
			}

			log.Print("Exiting and closing socket...")
			Start = false // Move server to down status
			ln.Close()
		}()

		//ln, err := net.Listen("tcp", ":8081")
		go a.listen_and_serve(ln, server_chan, &Start)

		time.Sleep(2)

		for {
			_, err = os.Stat(a.socket)
			if err != nil {
				break
			}
			time.Sleep(5)
		}

		if Start {
			log.Printf("Issue with the socket. %s. Closing it.", err)
			ln.Close()
			close(int_sig)
		}

		<-server_chan
		log.Print("http server is NOW off.")
		if !Start {
			os.Exit(0)
		}
		log.Print("Restarting the http server.")
		time.Sleep(2)
	}
}

func (a *__MYPLUGIN__App) listen_and_serve(ln net.Listener, server_chan chan bool, Start *bool) {
	log.Printf("httpd server: Starting service on socket '%s'", a.socket)

	router := NewRouter()

	srv := http.Server{Handler: router}
	err := srv.Serve(ln)
	if !*Start {
		log.Print("httpd server: Exiting...")
	} else {
		log.Printf("httpd server: Error detected: %s", err)
	}
	server_chan <- *Start
}

func (a *__MYPLUGIN__App) server_set() {
	if _, err := os.Stat(*a.params.socket_path); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(*a.params.socket_path, 0755)
		} else {
			kingpin.FatalIfError(err, "Unable to create '%s'\n", *a.params.socket_path)
		}
	}
	a.socket = path.Join(*a.params.socket_path, *a.params.socket_file)
}
