package server

import "net/http"

// Ping Healthcheck functionality for docker
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong\n"))
}
