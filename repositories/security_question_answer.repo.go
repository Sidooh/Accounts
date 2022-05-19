package repositories

import (
	Account "accounts.sidooh/models/account"
	SecurityQuestionAnswer "accounts.sidooh/models/security_question_answer"
	"accounts.sidooh/util"
	"errors"
	"strings"
)

func CheckAnswer(id uint, questionId uint, answer string) error {
	//	Get Account
	_, err := Account.ById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	//  Get Answer
	questionAnswer, err := SecurityQuestionAnswer.ByAccountAndQuestion(id, questionId)
	if err != nil {
		return errors.New("invalid question")
	}

	//	Check Answer

	if util.Compare(questionAnswer.Answer, strings.ToLower(answer)) {
		return nil
	} else {
		return errors.New("answer is incorrect")
	}
}

func HasSecurityQuestions(id uint) error {
	//	Get Account
	_, err := Account.ById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	//	Check Security Questions exists
	questionAnswers, err := SecurityQuestionAnswer.ByAccount(id)
	if err != nil {
		return errors.New("invalid account or questions")
	}

	if len(questionAnswers) == 3 {
		return nil
	}

	return errors.New("invalid answers found")
}
