package routes

import (
	"accounts.sidooh/models/invite"
	"accounts.sidooh/models/user"
	"accounts.sidooh/repositories"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
)

// TODO: Improve error handling, statuses, messages etc...
// TODO: Refactor, is very improper

func RegisterDashboardHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_DASHBOARD_URL+"/chart", func(context echo.Context) error {
		users, err := user.TimeSeriesCount(12)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		accounts, err := repositories.GetAccountsTimeData(12)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		invites, err := invite.TimeSeriesCount(12)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		var chart = struct {
			Users    interface{} `json:"users"`
			Accounts interface{} `json:"accounts"`
			Invites  interface{} `json:"invites"`
		}{
			users,
			accounts,
			invites,
		}

		return util.HandleSuccessResponse(context, chart)

	}, authMiddleware)

	e.GET(constants.API_DASHBOARD_URL+"/summaries", func(context echo.Context) error {
		users, err := user.Summaries()
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		accounts, err := repositories.GetAccountsSummary()
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		invites, err := invite.Summaries()
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		var summaries = struct {
			Users    interface{} `json:"users"`
			Accounts interface{} `json:"accounts"`
			Invites  interface{} `json:"invites"`
		}{
			users,
			accounts,
			invites,
		}

		return util.HandleSuccessResponse(context, summaries)
	}, authMiddleware)

	e.GET(constants.API_DASHBOARD_URL+"/recent-accounts", func(context echo.Context) error {
		data, err := repositories.GetAccounts(true, 20)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, data)
	}, authMiddleware)

	e.GET(constants.API_DASHBOARD_URL+"/recent-invites", func(context echo.Context) error {
		s, err := invite.All(20)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, s[:15])
	}, authMiddleware)
}
