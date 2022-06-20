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
	ErrBadRequest          = &ErrorResponse{Status: 400, Message: "Bad request."}
	ErrUnauthorized        = &ErrorResponse{Status: 401, Message: "Unauthorized."}
	ErrNotFound            = &ErrorResponse{Status: 404, Message: "Resource not found."}
	ErrNotAllowed          = &ErrorResponse{Status: 405, Message: "Not allowed."}
	ErrInternalServerError = &ErrorResponse{Status: 500, Message: "Internal Server Error."}
)

func (er *ErrorResponse) Render(w http.ResponseWriter, req *http.Request) error {
	render.Status(req, er.Status)
	return nil
}

func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:     err,
		Status:  400,
		Message: err.Error(),
	}
}
