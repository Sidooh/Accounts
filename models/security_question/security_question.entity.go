package security_question

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
)

type SecurityQuestion struct {
	models.Model

	Question string `json:"email"`
	Status   string `json:"status"`
}

func All() ([]SecurityQuestion, error) {
	conn := db.Connection()

	var questions []SecurityQuestion
	result := conn.Find(&questions)
	if result.Error != nil {
		return questions, result.Error
	}

	return questions, nil
}

//func CreateUser(u User) (User, error) {
//	conn := db.NewConnection().Conn
//	_, err := FindUserByEmail(u.Email)
//
//	if err == nil {
//		return User{}, errors.New("email is already taken")
//	}
//
//	u.Password, _ = util.ToHash(u.Password)
//
//	result := conn.Create(&u)
//	if result.Error != nil {
//		return User{}, errors.New("error creating user")
//	}
//
//	return u, nil
//}

//func AuthUser(u User) (User, error) {
//	user, err := FindUserByEmail(u.Email)
//
//	if err != nil {
//		return User{}, errors.New("invalid credentials")
//	}
//
//	res := util.Compare(user.Password, u.Password)
//
//	if !res {
//		return User{}, errors.New("invalid credentials")
//	}
//
//	return user, nil
//}

//func findAll(query interface{}, args interface{}) ([]SecurityQuestion, error) {
//	conn := db.Connection()
//
//	var questions []SecurityQuestion
//
//	result := conn.Where(query, args).Find(&questions)
//	if result.Error != nil {
//		return questions, result.Error
//	}
//
//	return questions, nil
//}
//
//func find(query interface{}, args interface{}) (SecurityQuestion, error) {
//	conn := db.Connection()
//
//	question := SecurityQuestion{}
//
//	result := conn.Where(query, args).First(&question)
//	if result.Error != nil {
//		return question, result.Error
//	}
//
//	return question, nil
//}
