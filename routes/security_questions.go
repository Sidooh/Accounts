package routes

import (
	"accounts.sidooh/middlewares"
	SecurityQuestion "accounts.sidooh/models/security_question"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
)

type CreateSecurityQuestionRequest struct {
	Question string `json:"question" validate:"required"`
}

func RegisterSecurityQuestionsHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/security-questions", func(context echo.Context) error {

		securityQuestions, err := SecurityQuestion.All()
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, securityQuestions)

	}, authMiddleware)

	e.POST(constants.API_URL+"/security-questions", func(context echo.Context) error {
		request := new(CreateSecurityQuestionRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		question, err := SecurityQuestion.CreateQuestion(SecurityQuestion.Model{
			Question: request.Question,
			Status:   "ACTIVE",
		})
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, question)

	}, authMiddleware)
}
