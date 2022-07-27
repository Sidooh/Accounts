package middlewares

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TokenAuth(secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(secret),
		TokenLookup: "header:Authorization",
		Claims:      &util.MyCustomClaims{},
		ErrorHandlerWithContext: func(err error, context echo.Context) error {
			unAuth := errors.NotAuthorizedError{Message: "Not Authorized"}
			return context.JSON(
				unAuth.Status(),
				unAuth.Errors(),
			)
		},
	})
}

// TODO: refresh token
//func RefreshTokenAuth(secret string) echo.MiddlewareFunc {
//	return middleware.JWTWithConfig(middleware.JWTConfig{
//		SigningKey:  []byte(secret),
//		TokenLookup: "cookie:jwt.sidooh,header:Authorization",
//		Claims:      &util.MyCustomClaims{},
//		ErrorHandlerWithContext: func(err error, context echo.Context) error {
//			unAuth := errors.NotAuthorizedError{Message: "Not Authorized"}
//			return context.JSON(
//				unAuth.Status(),
//				unAuth.Errors(),
//			)
//		},
//	})
//}
