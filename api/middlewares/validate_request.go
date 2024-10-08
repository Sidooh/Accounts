package middlewares

import (
	customErrors "accounts.sidooh/pkg"
	"accounts.sidooh/utils"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

// TODO: Should this be a pointer method?
func (cv CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {

		var validationErrors customErrors.RequestValidationErrors

		for _, err := range err.(validator.ValidationErrors) {

			msg := fmt.Sprintf("%v must be valid", err.Field())

			tag := err.Tag()
			if err.Param() != "" {
				tag += " " + err.Param()
			}
			validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, customErrors.ValidationError{
				Value:   err.Value().(string),
				Field:   err.Field(),
				Message: msg,
				Param:   tag,
			})
		}

		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(validationErrors.Status(), utils.ValidationErrorResponse(validationErrors.ValidationErrors))
	}
	return nil
}

func BindAndValidateRequest(context echo.Context, request interface{}) error {
	if err := context.Bind(request); err != nil {
		return err
	}

	if err := context.Validate(request); err != nil {
		return err
	}

	return nil
}

// TODO: Review validating request using generic
//func ValidateRequest[T any](context echo.Context, v T) error {
//	request := v
//	if err := BindAndValidateRequest(context, request); err != nil {
//		request.Id = context.Param("id")
//	}
//}
