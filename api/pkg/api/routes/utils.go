package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	hr "github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	l *zap.SugaredLogger
}

func (rw *responseWriter) WriteJSON(v interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	enc.SetIndent("", "  ")

	err := enc.Encode(v)
	if err != nil {
		rw.l.Errorf("Error writing JSON response: %v", err)
	}
	return err
}

type jsonError struct {
	Error string `json:"error"`
	Desc  string `json:"description,omitempty"`
	Code  int    `json:"code,omitempty"`
}

func jsonErrorFrom(err error, code int) jsonError {
	return jsonError{Error: err.Error(), Code: code}
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
