package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
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

func GenerateToken(user MyCustomClaims, validity time.Duration) (string, error) {
	MySigningKey := []byte(viper.GetString("JWT_KEY"))
	Iss := viper.GetString("TOKEN_ISSUER")
	Aud := viper.GetString("TOKEN_AUDIENCE")

	// Ensure validity is not more than a week and not less than 15mins. TODO: To review this
	if validity.Hours() > 7*24 || validity.Minutes() < 15 {
		validity = time.Duration(15) * time.Minute
	}

	claims := MyCustomClaims{
		user.Id,
		user.Email,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validity)),
			Subject:   strconv.Itoa(int(user.Id)),
			Issuer:    Iss,
			Audience:  jwt.ClaimStrings{Aud},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func SetTokenCookie(signedString string, ctx echo.Context) {
	env := strings.ToUpper(viper.GetString("APP_ENV"))

	if env != "PRODUCTION" {
		secureCookie = false
	}

	cookie := http.Cookie{
		Name:     "jwt.sidooh",
		Value:    signedString,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   secureCookie,
		HttpOnly: true,
		Path:     "/api",
	}
	ctx.SetCookie(&cookie)
}

func InvalidateToken(ctx echo.Context) {
	cookie := http.Cookie{
		Name:     "jwt.sidooh",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	ctx.SetCookie(&cookie)
}
