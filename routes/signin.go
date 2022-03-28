package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SignInRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=64"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func RegisterSignInHandler(e *echo.Echo) {
	e.POST("/api/users/signin", func(context echo.Context) error {

		request := new(SignInRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		user, err := User.AuthUser(User.User{
			Email:    request.Email,
			Password: strings.TrimSpace(request.Password),
		})
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		tokenData := util.MyCustomClaims{
			Id:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}
		token, _ := util.GenerateToken(tokenData)
		util.SetToken(token, context)

		return context.JSON(http.StatusOK, SignInResponse{Token: token})
	})
}
