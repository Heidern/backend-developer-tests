package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
)

type InvalidArgumentError struct {
	Err     error
	Message string
}

func (e *InvalidArgumentError) Error() string {
	return e.Message
}

type ResourceNotFoundError struct {
	Err     error
	Message string
}

func (e *ResourceNotFoundError) Error() string {
	return e.Message
}

type ApiErrorResponse struct {
	Error      error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status"`
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	apiError := ApiError(err)

	render.Status(r, apiError.StatusCode)
	render.JSON(w, r, apiError)

	logger := httplog.LogEntry(r.Context())

	if apiError.StatusCode == 500 {
		logger.Error().Msg(fmt.Sprintf("%v", err))
	} else {
		logger.Warn().Msg(fmt.Sprintf("%v", err))
	}
}

func ApiError(err error) *ApiErrorResponse {
	apiError := &ApiErrorResponse{
		Error:      err,
		StatusCode: 500,
		StatusText: err.Error(),
	}
	var invalidArgumentError *InvalidArgumentError

	if errors.As(err, &invalidArgumentError) {
		apiError.StatusCode = 400
	}

	var notFounderror *ResourceNotFoundError

	if errors.As(err, &notFounderror) {
		apiError.StatusCode = 404
	}

	if apiError.StatusCode == 500 {
		apiError.StatusText = "internal server error"
	}

	return apiError
}
