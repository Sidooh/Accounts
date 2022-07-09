package server

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	"accounts.sidooh/routes"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"
	"time"
)

func Setup() (*echo.Echo, string, *http2.Server) {
	fmt.Println("==== Starting Server ====")

	e := echo.New()
	e.HideBanner = true

	// Todo: Move to GetLogFile helper
	file := util.GetLogFile("server.log")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: file,
		Format: `{"time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","protocol":"${protocol}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Timeout())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true}))

	rateLimiterRequests := viper.GetFloat64("RATE_LIMIT")
	if rateLimiterRequests > 1 {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(rateLimiterRequests))))
	}

	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	authMiddlewareFunc := middlewares.TokenAuth(viper.GetString("JWT_KEY"))

	routes.RegisterCurrentUserHandler(e, authMiddlewareFunc)
	routes.RegisterSignInHandler(e)
	routes.RegisterSignUpHandler(e)
	routes.RegisterSignOutHandler(e, authMiddlewareFunc)

	// TODO: Review this properly
	//refreshAuthMiddlewareFunc := middlewares.RefreshTokenAuth(viper.GetString("JWT_KEY"))
	//routes.RegisterRefreshTokenHandler(e, refreshAuthMiddlewareFunc)

	// TODO: Add client auth endpoints clients/token & clients/refresh-token(token/refresh)?

	routes.RegisterAccountsHandler(e, authMiddlewareFunc)
	routes.RegisterInvitesHandler(e, authMiddlewareFunc)
	routes.RegisterUsersHandler(e, authMiddlewareFunc)
	routes.RegisterSecurityQuestionsHandler(e, authMiddlewareFunc)
	routes.RegisterSecurityQuestionAnswersHandler(e, authMiddlewareFunc)

	//-------------------
	// Custom middleware
	//-------------------
	// Stats
	statsMiddleware := middlewares.NewStats()
	e.Use(statsMiddleware.Process)
	e.Use(middlewares.ServerHeader)
	e.GET(constants.API_URL+"/stats", statsMiddleware.Handle) // Endpoint to get stats

	e.Any("*", func(context echo.Context) error {
		err := errors.NotFoundError{}

		return echo.NewHTTPError(err.Status(), err.Errors())
	})

	port := viper.GetString("PORT")
	if port == "" {
		port = "8000"
	}

	s := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	return e, port, s
}
