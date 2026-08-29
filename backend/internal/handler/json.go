package handler

import (
	"encoding/json"
	"net/http"
)

const maxRequestBodyBytes = 1 << 20

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return false
	}

	if decoder.More() {
		http.Error(w, "request body must contain one JSON value", http.StatusBadRequest)
		return false
	}

	return true
}