package security_question

import (
	"accounts.sidooh/models"
	"accounts.sidooh/pkg/db"
	"errors"
)

type Model struct {
	models.ModelID

	Question string `json:"question" gorm:"unique"`
	Status   string `json:"status"`

	models.ModelTimeStamps
}

func (Model) TableName() string {
	return "security_questions"
}

func All() ([]Model, error) {
	conn := db.Connection()

	var questions []Model
	result := conn.Order("id desc").Find(&questions)
	if result.Error != nil {
		return questions, result.Error
	}

	return questions, nil
}

func CreateQuestion(s Model) (Model, error) {
	conn := db.Connection()
	_, err := FindQuestion(s.Question)

	if err == nil {
		return Model{}, errors.New("question exists")
	}

	result := conn.Create(&s)
	if result.Error != nil {
		return Model{}, errors.New("error creating question")
	}

	return s, nil
}

func FindQuestion(question string) (Model, error) {
	return find("question = ?", question)
}

func find(query interface{}, args interface{}) (Model, error) {
	conn := db.Connection()

	secQn := Model{}

	result := conn.Where(query, args).First(&secQn)
	if result.Error != nil {
		return secQn, result.Error
	}

	return secQn, nil
}
