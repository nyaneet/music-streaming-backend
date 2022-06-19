package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Err     error  `json:"-"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var (
	ErrBadRequest = &ErrorResponse{Status: 400, Message: "Bad request."}
	ErrNotFound   = &ErrorResponse{Status: 404, Message: "Resource not found."}
	ErrNotAllowed = &ErrorResponse{Status: 405, Message: "Not allowed."}
)

func (er *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, er.Status)
	return nil
}

func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:     err,
		Status:  400,
		Message: err.Error(),
	}
}
