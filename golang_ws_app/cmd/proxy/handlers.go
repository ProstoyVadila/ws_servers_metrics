package main

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("go_ws_app"))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
