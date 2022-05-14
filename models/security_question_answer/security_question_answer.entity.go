package security_question_answer

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	Account "accounts.sidooh/models/account"
	SecurityQuestion "accounts.sidooh/models/security_question"
	"errors"
)

type Model struct {
	models.ModelID

	Answer string

	QuestionID uint `json:"-"`
	AccountID  uint `json:"-"`

	models.ModelTimeStamps
}

type ModelWithQuestion struct {
	Model

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

func ByAccountIdWithQuestion(id uint) ([]ModelWithQuestion, error) {
	var questions []ModelWithQuestion

	result := db.Connection().Joins("Question").Find(&questions, id)
	if result.Error != nil {
		return questions, result.Error
	}

	return questions, nil
}

// TODO: Check whether using pointers here saves memory
func Create(s Model) (Model, error) {
	_, err := ByAccountAndQuestion(s.AccountID, s.QuestionID)
	if err == nil {
		return Model{}, errors.New("question and answer already exists")
	}

	result := db.Connection().Create(&s)
	if result.Error != nil {
		return Model{}, errors.New("error creating answer")
	}

	return s, nil
}

func ByAccountAndQuestion(accountId uint, questionId uint) (Model, error) {
	model := Model{}

	result := db.Connection().Where("account_id = ? and question_id = ?", accountId, questionId).First(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func find(query interface{}, args interface{}) (Model, error) {
	model := Model{}

	result := db.Connection().Where(query, args).First(&model)
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}
