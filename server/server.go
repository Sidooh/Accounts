package server

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	"accounts.sidooh/routes"
	"accounts.sidooh/util"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"os"
	"path/filepath"
	"time"
)

func Setup() (*echo.Echo, string, *http2.Server) {
	e := echo.New()
	e.HideBanner = true

	// Todo: Move to GetLogFile helper
	pwd, err := os.Getwd()
	file := util.GetFile(filepath.Join(pwd, "/logs/", "server.log"))
	if err != nil || file == nil {
		// Handle error
		panic("could not open file")
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: file,
		Format: `{"time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","protocol":"${protocol}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Secure())
	e.Use(middleware.Timeout())

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

	port := viper.GetString("PORT")
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
