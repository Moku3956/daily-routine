package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Status  int      `json:"status"`
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string, details ...string) {
	errText := map[int]string{
		http.StatusBadRequest:          "Bad Request",
		http.StatusUnauthorized:        "Unauthorized",
		http.StatusForbidden:           "Forbidden",
		http.StatusNotFound:            "Not Found",
		http.StatusConflict:            "Conflict",
		http.StatusUnprocessableEntity: "Unprocessable Entity",
	}[status]
	if errText == "" {
		errText = "Internal Server Error"
	}
	resp := errorResponse{Status: status, Error: errText, Message: message}
	if len(details) > 0 {
		resp.Details = details
	}
	writeJSON(w, status, resp)
}
