package security_question_answers

import (
	"accounts.sidooh/pkg/repositories/accounts"
	"accounts.sidooh/utils"
	"errors"
	"strings"
)

func CheckAnswer(id uint, questionId uint, answer string) error {
	//	Get Account
	_, err := accounts.ReadById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	//  Get Answer
	questionAnswer, err := ReadByAccountAndQuestion(id, questionId)
	if err != nil {
		return errors.New("invalid question")
	}

	//	Check Answer

	if utils.Compare(questionAnswer.Answer, strings.ToLower(answer)) {
		return nil
	} else {
		return errors.New("answer is incorrect")
	}
}

func HasSecurityQuestions(id uint) error {
	//	Get Account
	_, err := accounts.ReadById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	//	Check Security Questions exists
	questionAnswers, err := ReadByAccountId(id)
	if err != nil {
		return errors.New("invalid account or questions")
	}

	if len(questionAnswers) == 3 {
		return nil
	}

	return errors.New("invalid answers found")
}
