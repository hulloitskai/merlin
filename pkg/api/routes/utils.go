package routes

import (
	"encoding/json"
	"net/http"
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
