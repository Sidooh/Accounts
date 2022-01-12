package server

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	"accounts.sidooh/routes"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/http2"
	"os"
	"time"
)

func Setup() (*echo.Echo, string, *http2.Server) {
	e := echo.New()
	e.HideBanner = true

	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	routes.RegisterCurrentUserHandler(e)
	routes.RegisterSignInHandler(e)
	routes.RegisterSignUpHandler(e)
	routes.RegisterSignOutHandler(e)

	routes.RegisterAccountsHandler(e)
	routes.RegisterReferralsHandler(e)

	e.Any("*", func(context echo.Context) error {
		err := errors.NotFoundError{}

		return echo.NewHTTPError(err.Status(), err.Errors())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	s := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	return e, port, s
}
