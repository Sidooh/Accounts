package routes

import (
	"accounts.sidooh/api/middlewares"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/pkg"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
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
			return echo.NewHTTPError(400, pkg.BadRequestError{Message: err.Error()}.Errors())
		}

		tokenData := utils.MyCustomClaims{
			Id:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}
		validity := time.Duration(viper.GetInt("ACCESS_TOKEN_VALIDITY")) * time.Minute
		accessToken, _ := utils.GenerateToken(tokenData, validity)

		return context.JSON(http.StatusOK, SignInResponse{AccessToken: accessToken})
	})
}
