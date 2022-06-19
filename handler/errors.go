package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

var (
	ErrBadRequest = &ErrorResponse{StatusCode: 400, Message: "Bad request."}
	ErrNotFound   = &ErrorResponse{StatusCode: 404, Message: "Resource not found."}
	ErrNotAllowed = &ErrorResponse{StatusCode: 405, Message: "Not allowed."}
)

func (er *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, er.StatusCode)
	return nil
}

func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}
