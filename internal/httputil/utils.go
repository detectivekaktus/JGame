package httputil

import (
	"encoding/json"
	"fmt"
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

// Asserts the Content-Type header value is equal to `expected` value.
// If the Content-Type does not match the `expected` value, false is returned
// and the error is sent via HTTP to the client. Otherwise, true is returned
// and no actions are done.
func IsContentType(w http.ResponseWriter, r *http.Request, expected string) bool {
	if ctnType := r.Header.Get("Content-Type"); ctnType != expected {
		SendErrorMessage(w, http.StatusBadRequest, "Content error",
			fmt.Sprintf("Excepted to find Content-Type `%s`, found `%s`.", expected, ctnType))
		return false
	}
	return true
}

// Returns true of the request has provided body, false otherwise.
func HasContent(r *http.Request) bool {
	return r.ContentLength > 0
}
