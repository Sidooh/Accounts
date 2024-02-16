package routes

import (
	accountsRepo "accounts.sidooh/pkg/repositories/accounts"
	invitesRepo "accounts.sidooh/pkg/repositories/invites"
	usersRepo "accounts.sidooh/pkg/repositories/users"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
)

func RegisterAnalyticsHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_ANALYTICS_URL+"/accounts", func(ctx echo.Context) error {
		dataset, err := accountsRepo.GetAccountsTimeSeries()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, dataset)
	}, authMiddleware)

	e.GET(constants.API_ANALYTICS_URL+"/users", func(ctx echo.Context) error {
		dataset, err := usersRepo.ReadTimeSeriesCount()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, dataset)
	}, authMiddleware)

	e.GET(constants.API_ANALYTICS_URL+"/invites", func(ctx echo.Context) error {
		dataset, err := invitesRepo.ReadTimeSeriesCount()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, dataset)
	}, authMiddleware)
}
