package util

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var secureCookie = true

type MyCustomClaims struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateToken(user MyCustomClaims) (string, error) {
	MySigningKey := []byte(viper.GetString("JWT_KEY"))

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
