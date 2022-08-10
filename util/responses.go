package util

import (
	customErrors "accounts.sidooh/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type JsonResponse struct {
	Result  int         `json:"result"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(data interface{}) JsonResponse {
	return JsonResponse{
		Result: 1,
		Data:   data,
	}
}

func ErrorResponse(message string, errors interface{}) JsonResponse {
	return JsonResponse{
		Result:  0,
		Message: message,
		Errors:  errors,
	}
}

func ServerErrorResponse() JsonResponse {
	return JsonResponse{
		Result:  0,
		Message: "Something went wrong, please try again.",
	}
}

func NotFoundErrorResponse() JsonResponse {
	return JsonResponse{
		Result:  0,
		Message: "Not Found",
	}
}

func UnauthenticatedErrorResponse() JsonResponse {
	return JsonResponse{
		Result:  0,
		Message: "Unauthenticated",
	}
}

func ValidationErrorResponse(errors interface{}) JsonResponse {
	return ErrorResponse("The request is invalid", errors)
}

func IdValidationErrorResponse(value string) JsonResponse {
	err := customErrors.ValidationError{
		Value:   value,
		Field:   "Id",
		Message: "Id must be valid",
		Param:   "must be valid numeric",
	}

	return ValidationErrorResponse(err)
}

func PhoneValidationErrorResponse(value string) JsonResponse {
	err := customErrors.ValidationError{
		Value:   value,
		Field:   "Phone",
		Message: "Phone must be valid",
		Param:   "must be valid phone number",
	}

	return ValidationErrorResponse(err)
}

func HandleErrorResponse(ctx echo.Context, err error) error {
	fmt.Println(err)

	if err.Error() == "record not found" {
		return ctx.JSON(http.StatusNotFound, NotFoundErrorResponse())
	}

	return ctx.JSON(http.StatusInternalServerError, ServerErrorResponse())
}

func HandleSuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, SuccessResponse(data))
}
