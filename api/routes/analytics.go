package routes

import (
	accountsRepo "accounts.sidooh/pkg/repositories/accounts"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
)

func RegisterAnalyticsHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_DASHBOARD_URL+"/analytics/accounts", func(ctx echo.Context) error {
		dataset, err := accountsRepo.GetAccountsTimeSeries()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, dataset)
	}, authMiddleware)
}
