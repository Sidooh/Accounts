package pkg

import (
	"net/http"
)

type ValidationError struct {
	Value   string `json:"value"`
	Field   string `json:"field"`
	Message string `json:"message"`
	Param   string `json:"param"`
}

type RequestValidationErrors struct {
	ValidationErrors []ValidationError
}

type DatabaseConnectionError struct {
	Reason string `json:"message"`
}

type NotFoundError struct{}

type BadRequestError struct {
	Message string `json:"message"`
}

type NotAuthorizedError struct {
	Message string `json:"message"`
}

func (d NotFoundError) Error() string {
	panic("Not Found Error")
}

func (d ValidationError) Error() string {
	panic("Validation Error")
}

func (d DatabaseConnectionError) Error() string {
	panic("DB Conn Error")
}

func (d BadRequestError) Error() string {
	panic("Bad Request Error")
}

func (d NotAuthorizedError) Error() string {
	panic("Not Authorized Error")
}

func (d RequestValidationErrors) Status() int {
	return http.StatusUnprocessableEntity
}

func (d DatabaseConnectionError) Status() int {
	return http.StatusInternalServerError
}

func (d NotFoundError) Status() int {
	return http.StatusNotFound
}

func (d BadRequestError) Status() int {
	return http.StatusBadRequest
}

func (d NotAuthorizedError) Status() int {
	return http.StatusUnauthorized
}

func (d DatabaseConnectionError) Errors() map[string][]error {
	return map[string][]error{
		"errors": {d},
	}
}

func (d RequestValidationErrors) Errors() map[string][]ValidationError {
	return map[string][]ValidationError{
		"errors": d.ValidationErrors,
	}
}

func (d NotFoundError) Errors() map[string][]string {
	return map[string][]string{
		"errors": {"Not Found"},
	}
}

func (d BadRequestError) Errors() map[string][]error {
	return map[string][]error{
		"errors": {d},
	}
}

func (d NotAuthorizedError) Errors() map[string][]error {
	return map[string][]error{
		"errors": {d},
	}
}
