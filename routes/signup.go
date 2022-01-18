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

type SignUpRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=4,max=20"`
	Username string `json:"username" form:"username"`
}

func RegisterSignUpHandler(e *echo.Echo) {
	go e.POST("/api/users/signup", func(context echo.Context) error {
		return handlerFunc(context)
	})
}

func handlerFunc(context echo.Context) error {
	request := new(SignUpRequest)
	if err := middlewares.BindAndValidateRequest(context, request); err != nil {
		return err
	}

	username := request.Username
	if username == "" {
		username = request.Email
	}

	user, err := User.CreateUser(User.User{
		Email:    request.Email,
		Password: strings.TrimSpace(request.Password),
		Username: username,
	})
	if err != nil {
		return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
	}

	tokenData := util.MyCustomClaims{
		Id:    user.ID,
		Email: user.Email,
	}
	token, _ := util.GenerateToken(tokenData)
	util.SetToken(token, context)

	return context.JSON(http.StatusCreated, user)
}
