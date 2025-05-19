package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorMsg struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Sets `status` status code for the `w` writer.
// Writes a JSON error message with the `status` status, `error`
// title, and `msg` message.
//
// You must first set the Content-Type header to application/json
// before calling this function.
func SendErrorMessage(w http.ResponseWriter, status int, error, msg string) {
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorMsg{
		Error: error,
		Message: msg,
	})
}
