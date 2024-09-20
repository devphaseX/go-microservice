package main

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest      = errors.New("Bad Request")
	ErrInternalFailure = errors.New("Internal failure")
	ErrNotFound        = errors.New("Not found")
)

type Error struct {
	appErr error
	srvErr error
}

func (e Error) AppErr() error {
	return e.appErr
}

func (e Error) SrvErr() error {
	return e.srvErr
}

func NewError(appErr error, srvErr error) error {
	return Error{
		appErr: appErr,
		srvErr: srvErr,
	}
}

func (e Error) Error() string {
	return errors.Join(e.srvErr, e.appErr).Error()
}

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

func FromError(err error) ApiError {
	var (
		apiError ApiError
		svcError Error
	)

	if errors.As(err, &svcError) {
		apiError.Message = svcError.AppErr().Error()
		svcErr := svcError.SrvErr()

		switch svcErr {
		case ErrBadRequest:
			apiError.Status = http.StatusBadGateway

		case ErrInternalFailure:
			apiError.Status = http.StatusInternalServerError

		case ErrNotFound:
			apiError.Status = http.StatusNotFound

		default:
			apiError.Status = http.StatusInternalServerError
		}

	}
	return apiError
}
