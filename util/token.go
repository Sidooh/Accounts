package util

import (
	"accounts.sidooh/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var MySigningKey = []byte(viper.GetString("JWT_KEY"))
var secureCookie = true

type MyCustomClaims struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

var CustomJWTMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey:  MySigningKey,
	TokenLookup: "cookie:jwt,header:Authorization",
	Claims:      &MyCustomClaims{},
	ErrorHandlerWithContext: func(err error, context echo.Context) error {
		unAuth := errors.NotAuthorizedError{Message: "Not Authorized"}
		return context.JSON(
			unAuth.Status(),
			unAuth.Errors(),
		)
	},
})

func GenerateToken(user MyCustomClaims) (string, error) {
	claims := MyCustomClaims{
		user.Id,
		user.Email,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func SetToken(signedString string, ctx echo.Context) {
	env := strings.ToUpper(viper.GetString("APP_ENV"))

	if env != "PRODUCTION" {
		secureCookie = false
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedString,
		Expires:  time.Now().Add(15 * time.Minute),
		Secure:   secureCookie,
		HttpOnly: true,
		Path:     "/api",
	}
	ctx.SetCookie(&cookie)
}

func InvalidateToken(ctx echo.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	ctx.SetCookie(&cookie)
}
