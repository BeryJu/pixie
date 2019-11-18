package base

import "net/http"

// ServingFile http.File with an additional Method to serve
type ServingFile interface {
	http.File
	Serve(w http.ResponseWriter, r *http.Request)
}
