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

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
func BadRequest(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}
func ServerError(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 500,
		StatusText: "Internal server error",
		Message:    err.Error(),
	}
}
