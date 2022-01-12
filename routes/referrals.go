package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	Referral "accounts.sidooh/models/referral"
	"accounts.sidooh/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateReferralRequest struct {
	Phone        string `json:"phone" form:"phone" validate:"required,numeric"`
	RefereePhone string `json:"referee_phone" form:"referee_phone" validate:"required,numeric"`
}

func RegisterReferralsHandler(e *echo.Echo) {
	e.GET("/api/referrals", func(context echo.Context) error {

		referrals, err := Referral.All()
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, referrals)

	}, services.CustomJWTMiddleware)

	e.GET("/api/referrals/:phone", func(context echo.Context) error {

		phone := context.Param("phone")

		phone, err := services.GetPhoneByCountry("KE", phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		referral, err := Referral.ByPhone(nil, phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, referral)

	}, services.CustomJWTMiddleware)

	e.POST("/api/referrals", func(context echo.Context) error {

		request := new(CreateReferralRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := services.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		referral, err := Referral.Create(Referral.Model{
			RefereePhone: phone,
		}, phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, referral)
	})
}
