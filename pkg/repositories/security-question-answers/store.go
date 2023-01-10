package security_question_answers

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/utils"
	"errors"
	"strings"
)

// TODO: Check whether using pointers here saves memory

func Create(s entities.QuestionAnswer) (entities.QuestionAnswer, error) {
	_, err := ReadByAccountAndQuestion(s.AccountID, s.QuestionID)
	if err == nil {
		return entities.QuestionAnswer{}, errors.New("question and answer already exists")
	}

	s.Answer, _ = utils.ToHash(strings.ToLower(s.Answer))

	result := db.Connection().Create(&s)
	if result.Error != nil {
		return entities.QuestionAnswer{}, errors.New("error creating answer")
	}

	return s, nil
}

func ReadByAccountId(accountId uint) ([]entities.QuestionAnswer, error) {
	var model []entities.QuestionAnswer

	result := db.Connection().Where("account_id = ?", accountId).Find(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func ReadByAccountAndQuestion(accountId uint, questionId uint) (entities.QuestionAnswer, error) {
	model := entities.QuestionAnswer{}

	result := db.Connection().Where("account_id = ? and question_id = ?", accountId, questionId).First(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func ReadByAccountIdWithQuestion(id uint) ([]entities.QuestionAnswerWithQuestion, error) {
	var questions []entities.QuestionAnswerWithQuestion

	// TODO: Should we check for existence of account first?
	// TODO: Can we flatten the results from here? (like array flatten)
	result := db.Connection().Where("account_id = ?", id).Joins("Question").Find(&questions)
	if result.Error != nil {
		return questions, result.Error
	}

	return questions, nil
}

func find(query interface{}, args interface{}) (entities.QuestionAnswer, error) {
	model := entities.QuestionAnswer{}

	result := db.Connection().Where(query, args).First(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}
