package routes

import (
	"accounts.sidooh/api/middlewares"
	"accounts.sidooh/pkg/repositories/users"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type SearchEmailRequest struct {
	Email string `query:"email" validate:"required,min=2,max=32"`
}

type UserByIdRequest struct {
	Id          string `param:"id" validate:"required,numeric,number,min=1"`
	WithAccount string `query:"with_account" validate:"omitempty,oneof=true false"`
}

func RegisterUsersHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/users", func(context echo.Context) error {
		fetchedUsers, err := users.ReadAll()
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, fetchedUsers)
	}, authMiddleware)

	e.GET(constants.API_URL+"/users/:id", func(context echo.Context) error {
		request := new(UserByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		//if request.WithAccount == "true" {
		//	//TODO: To implement
		//} else {
		//	user, err := User.FindUserById(uint(id))
		//	if err != nil {
		//		return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		//	}
		//
		//	return context.JSON(http.StatusOK, user)
		//}

		user, err := users.ReadById(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, user)
	}, authMiddleware)

	e.GET(constants.API_URL+"/users/search", func(context echo.Context) error {
		request := new(SearchEmailRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		fetchedUsers, err := users.SearchByEmail(request.Email)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, fetchedUsers)
	}, authMiddleware)

	e.POST(constants.API_URL+"/users/:id/reset-password", func(ctx echo.Context) error {
		request := new(UserByIdRequest)
		if err := middlewares.BindAndValidateRequest(ctx, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		if err = users.ResetPassword(uint(id)); err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, true)
	})
}
