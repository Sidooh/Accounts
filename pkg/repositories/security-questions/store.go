package security_questions

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"errors"
)

func ReadAll() ([]entities.Question, error) {
	conn := db.Connection()

	var questions []entities.Question
	result := conn.Order("id desc").Find(&questions)
	if result.Error != nil {
		return questions, result.Error
	}

	return questions, nil
}

func Create(s entities.Question) (entities.Question, error) {
	conn := db.Connection()
	if _, err := ReadByQuestion(s.Question); err == nil {
		return entities.Question{}, errors.New("question exists")
	}

	result := conn.Create(&s)
	if result.Error != nil {
		return entities.Question{}, errors.New("error creating question")
	}

	return s, nil
}

func ReadByQuestion(question string) (entities.Question, error) {
	return findQuestion("question = ?", question)
}

func findQuestion(query interface{}, args interface{}) (entities.Question, error) {
	conn := db.Connection()

	secQn := entities.Question{}

	result := conn.Where(query, args).First(&secQn)
	if result.Error != nil {
		return secQn, result.Error
	}

	return secQn, nil
}
