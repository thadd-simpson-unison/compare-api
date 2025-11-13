package Jsend

import (
	"encoding/json"
	"net/http"
)

// Status constants
const (
	StatusError    = "error"
	StatusFail     = "fail"
	StatusRedirect = "redirect"
	StatusSuccess  = "success"
)

// Body contains
type Body struct {
	// The status indicates the execution result of request,
	// it can be one of "success", "fail" and "error".
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// CONSTRUCTORS

// New returns a success body with the given data.
func New(data interface{}) Body {
	return Body{
		Status: StatusSuccess,
		Data:   data,
	}
}

// NewFail returns a fail body with the given data.
func NewFail(data interface{}) Body {
	return Body{
		Status: StatusFail,
		Data:   data,
	}
}

// NewError returns a error body with given message.
func NewError(message string, code int, data interface{}) Body {
	return Body{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Data:    data,
	}
}

// Write writes the body to http.ResponseWriter.
// If necessary, the status code can be specified through the third parameter.
func Write(w http.ResponseWriter, body Body, statuses ...int) error {
	w.Header().Set("Content-Type", "application/json")

	if len(statuses) > 0 {
		w.WriteHeader(statuses[0])
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// Pass plain text back to the client without json formatting
func WritePlain(w http.ResponseWriter, body string, statuses ...int) {
	w.Header().Set("Content-Type", "text/plain")

	if len(statuses) > 0 {
		w.WriteHeader(statuses[0])
	}

	b := []byte(body)

	w.Write(b)
}

// JSENDERS

// Error writes error body with the given message.
func Error(w http.ResponseWriter, message string, statuses ...int) error {
	return Write(w, NewError(message, 0, nil), statuses...)
}

// ErrorCode writes error body with the given message and code.
func ErrorCode(w http.ResponseWriter, message string, code int, statuses ...int) error {
	return Write(w, NewError(message, code, nil), statuses...)
}

// ErrorCodeData writes error body with the given message, code and data.
func ErrorCodeData(w http.ResponseWriter, message string, code int, data interface{}, statuses ...int) error {
	return Write(w, NewError(message, code, data), statuses...)
}

// Fail writes failed body with the given data.
func Fail(w http.ResponseWriter, data interface{}, statuses ...int) error {
	return Write(w, NewFail(data), statuses...)
}

// Redirect a string url back to the client
func Redirect(w http.ResponseWriter, url string, statuses ...int) {
	WritePlain(w, url, statuses...)
}

// Success writes successful body with the given data.
func Success(w http.ResponseWriter, data interface{}, statuses ...int) error {
	return Write(w, New(data), statuses...)
}
