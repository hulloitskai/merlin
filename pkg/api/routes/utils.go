package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	hr "github.com/julienschmidt/httprouter"
)

type responseWriter struct {
	http.ResponseWriter
}

func (rw *responseWriter) WriteJSON(v interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

type jsonError struct {
	Error string `json:"error"`
}

func jsonErrorFrom(err error) jsonError {
	return jsonError{Error: err.Error()}
}

func handleTrailingSlashRedir(w http.ResponseWriter, r *http.Request,
	_ hr.Params) {
	var (
		prefix = r.Header.Get("X-Forwarded-Prefix")
		path   = strings.TrimSuffix(prefix, "/") + r.URL.Path
	)

	if len(path) > 1 && (path[len(path)-1] == '/') {
		r.URL.Path = path[:len(path)-1]
	} else {
		r.URL.Path = path + "/"
	}

	code := 301 // Permanent redirect, request with GET method
	if r.Method != "GET" {
		// Temporary redirect, request with same method
		// As of Go 1.3, Go does not support status code 308.
		code = 307
	}
	http.Redirect(w, r, r.URL.String(), code)
}
