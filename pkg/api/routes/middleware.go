package routes

import (
	"net/http"
)

func corsMiddleware(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		hf(w, r)
	}
}
