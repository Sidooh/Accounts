package security_question_answer

import (
	"accounts.sidooh/models"
	Account "accounts.sidooh/models/account"
	SecurityQuestion "accounts.sidooh/models/security_question"
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/utils"
	"errors"
	"strings"
)

type Model struct {
	models.ModelID

	Answer string `json:"-"`

	QuestionID uint `json:"-"`
	AccountID  uint `json:"-"`

	models.ModelTimeStamps
}

type ModelWithQuestion struct {
	// TODO: Can we flatten the results from here? (like array flatten),
	// 	related: Since this only brings back id, is it necessary?
	Model /*`json:"-"`*/

	Question SecurityQuestion.Model `json:"question"`
}

type ModelWithAccountAndQuestion struct {
	ModelWithQuestion

	Account Account.Model `json:"account"`
}

func (Model) TableName() string {
	return "security_question_answers"
}

func (ModelWithQuestion) TableName() string {
	return "security_question_answers"
}

func (ModelWithAccountAndQuestion) TableName() string {
	return "security_question_answers"
}

// TODO: Check whether using pointers here saves memory
func Create(s Model) (Model, error) {
	_, err := ByAccountAndQuestion(s.AccountID, s.QuestionID)
	if err == nil {
		return Model{}, errors.New("question and answer already exists")
	}

	s.Answer, _ = utils.ToHash(strings.ToLower(s.Answer))

	result := db.Connection().Create(&s)
	if result.Error != nil {
		return Model{}, errors.New("error creating answer")
	}

	return s, nil
}

func ByAccount(accountId uint) ([]Model, error) {
	var model []Model

	result := db.Connection().Where("account_id = ?", accountId).Find(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func ByAccountAndQuestion(accountId uint, questionId uint) (Model, error) {
	model := Model{}

	result := db.Connection().Where("account_id = ? and question_id = ?", accountId, questionId).First(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func ByAccountIdWithQuestion(id uint) ([]ModelWithQuestion, error) {
	var questions []ModelWithQuestion

	// TODO: Should we check for existence of account first?
	// TODO: Can we flatten the results from here? (like array flatten)
	result := db.Connection().Where("account_id = ?", id).Joins("Question").Find(&questions)
	if result.Error != nil {
		return questions, result.Error
	}

	return questions, nil
}

func find(query interface{}, args interface{}) (Model, error) {
	model := Model{}

	result := db.Connection().Where(query, args).First(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}
