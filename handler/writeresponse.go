package handler

import (
	"encoding/json"
	"net/http"
)

func writeErrResponse(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// ValidResponse writes a successful JSON response
func writeValidResponse(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
