package main

const template_rest_log = `package main

// This file has been created by "go generate" as initial code. go generate will never update it, except if you remove it.

// So, update it for your need.

import (
    "log"
    "net/http"
    "time"
)

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w, r)

        log.Printf(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.RequestURI,
            name,
            time.Since(start),
        )
    })
}`
