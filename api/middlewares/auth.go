package middlewares

import (
	"accounts.sidooh/pkg"
	"accounts.sidooh/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TokenAuth(secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(secret),
		TokenLookup: "header:Authorization",
		Claims:      &utils.MyCustomClaims{},
		ErrorHandlerWithContext: func(err error, context echo.Context) error {
			unAuth := pkg.NotAuthorizedError{Message: "Not Authorized"}
			//return context.JSON(
			//	unAuth.Status(),
			//	unAuth.Errors(),
			//)

			return context.JSON(unAuth.Status(), utils.UnauthenticatedErrorResponse())
		},
	})
}

// TODO: refresh token
//func RefreshTokenAuth(secret string) echo.MiddlewareFunc {
//	return middleware.JWTWithConfig(middleware.JWTConfig{
//		SigningKey:  []byte(secret),
//		TokenLookup: "cookie:jwt.sidooh,header:Authorization",
//		Claims:      &utils.MyCustomClaims{},
//		ErrorHandlerWithContext: func(err error, context echo.Context) error {
//			unAuth := errors.NotAuthorizedError{Message: "Not Authorized"}
//			return context.JSON(
//				unAuth.Status(),
//				unAuth.Errors(),
//			)
//		},
//	})
//}
