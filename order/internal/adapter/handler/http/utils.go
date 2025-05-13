package http

import (
	"encoding/json"
	"net/http"
)

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(w, status, envelope{message})
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	// app.logger.Error("Bad request:", zap.String("error", err.Error()), zap.String("path", r.URL.Path), zap.String("method", r.Method))
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{data})
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}
