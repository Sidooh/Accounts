package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	SecurityQuestionAnswer "accounts.sidooh/models/security_question_answer"
	"accounts.sidooh/repositories"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

//type QuestionAnswerRequest struct {
//	QuestionId string `json:"question_id" validate:"required"`
//	Answer     string `json:"answer" validate:"required"`
//}

type CreateSecurityQuestionAnswersRequest struct {
	Id string `param:"id" validate:"required,numeric,min=1"`
	//TODO: Can we have bulk creation?
	//Questions []QuestionAnswerRequest `json:"questions" validate:"required"`

	QuestionId string `json:"question_id" validate:"required"`
	Answer     string `json:"answer" validate:"required"`
}

type QuestionsByAccountIdRequest struct {
	Id            string `param:"id" validate:"required,numeric,min=1"`
	WithQuestions string `query:"with_user" validate:"omitempty,oneof=true false"`
}

func RegisterSecurityQuestionAnswersHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/accounts/:id/security-questions", func(context echo.Context) error {
		request := new(QuestionsByAccountIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			request.Id = context.Param("id")
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		securityQuestions, err := SecurityQuestionAnswer.ByAccountIdWithQuestion(uint(id))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, securityQuestions)

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/security-questions/answers", func(context echo.Context) error {
		request := new(CreateSecurityQuestionAnswersRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			request.Id = context.Param("id")
		}

		id, err := strconv.Atoi(request.Id)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		questionId, err := strconv.Atoi(request.QuestionId)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		question, err := SecurityQuestionAnswer.Create(SecurityQuestionAnswer.Model{
			AccountID:  uint(id),
			QuestionID: uint(questionId),
			Answer:     request.Answer,
		})
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, question)

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/security-questions/check", func(context echo.Context) error {
		request := new(CreateSecurityQuestionAnswersRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			request.Id = context.Param("id")
		}

		id, err := strconv.Atoi(request.Id)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		questionId, err := strconv.Atoi(request.QuestionId)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.CheckAnswer(uint(id), uint(questionId), request.Answer)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/has-security-questions", func(context echo.Context) error {

		request := new(QuestionsByAccountIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.HasSecurityQuestions(uint(id))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]bool{
			"message": true,
		})

	}, authMiddleware)
}
