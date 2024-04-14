package routes

import (
	"accounts.sidooh/api/middlewares"
	"accounts.sidooh/pkg/repositories/otp"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OtpRequest struct {
	Phone string `json:"phone" validate:"required,numeric,min=9"`
	Otp   int    `json:"otp,omitempty,numeric,min=100000,max=999999"`
}

func RegisterOtpHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.POST(constants.API_URL+"/otp/generate", func(context echo.Context) error {
		request := new(OtpRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := utils.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.PhoneValidationErrorResponse(request.Phone))
		}

		err = otp.GenerateOTP(phone)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, false)
		}

		return context.JSON(http.StatusOK, true)
	}, authMiddleware)

	e.POST(constants.API_URL+"/otp/verify", func(context echo.Context) error {
		request := new(OtpRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := utils.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.PhoneValidationErrorResponse(request.Phone))
		}

		err = otp.ValidateOTP(phone, request.Otp)
		if err != nil {
			return context.JSON(http.StatusBadRequest, false)
		}

		return context.JSON(http.StatusOK, true)
	}, authMiddleware)
}
