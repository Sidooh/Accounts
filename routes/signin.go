package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

type SignInRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=64"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func RegisterSignInHandler(e *echo.Echo) {
	e.POST(constants.API_URL+"/users/signin", func(context echo.Context) error {

		request := new(SignInRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		user, err := User.AuthUser(User.Model{
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
		validity := time.Duration(viper.GetInt("ACCESS_TOKEN_VALIDITY")) * time.Minute
		accessToken, _ := util.GenerateToken(tokenData, validity)

		//validity = time.Duration(viper.GetInt("REFRESH_TOKEN_VALIDITY"))
		//refreshToken, _ := util.GenerateToken(util.MyCustomClaims{Id: user.ID}, validity)

		//util.SetTokenCookie(refreshToken, context)

		return context.JSON(http.StatusOK, SignInResponse{AccessToken: accessToken})
	})
}
