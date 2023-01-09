package routes

import (
	"accounts.sidooh/api/middlewares"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/pkg"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SignUpRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=64"`
	Username string `json:"username" form:"username"`
}

func RegisterSignUpHandler(e *echo.Echo) {
	e.POST(constants.API_URL+"/users/signup", func(context echo.Context) error {
		return signup(context)
	})
}

func signup(context echo.Context) error {
	request := new(SignUpRequest)
	if err := middlewares.BindAndValidateRequest(context, request); err != nil {
		return err
	}

	username := request.Username
	if username == "" {
		username = request.Email
	}

	user, err := User.CreateUser(User.Model{
		Email:    request.Email,
		Password: strings.TrimSpace(request.Password),
		Username: username,
	})
	if err != nil {
		return echo.NewHTTPError(400, pkg.BadRequestError{Message: err.Error()}.Errors())
	}

	// TODO: should user have to login after signing up?, I am thinking yes!
	//tokenData := utils.MyCustomClaims{
	//	Id:    user.ID,
	//	Email: user.Email,
	//}
	//token, _ := utils.GenerateToken(tokenData)
	//utils.SetToken(token, context)

	return context.JSON(http.StatusCreated, user)
}
