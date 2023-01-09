package routes

import (
	"accounts.sidooh/api/middlewares"
	SecurityQuestion "accounts.sidooh/models/security_question"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
)

type CreateSecurityQuestionRequest struct {
	Question string `json:"question" validate:"required"`
}

func RegisterSecurityQuestionsHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/security-questions", func(context echo.Context) error {

		securityQuestions, err := SecurityQuestion.All()
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, securityQuestions)

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
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, question)

	}, authMiddleware)
}
