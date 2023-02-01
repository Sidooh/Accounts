package users

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/utils"
	"errors"
	"fmt"
	"time"
)

func ReadAll() ([]entities.User, error) {
	conn := db.Connection()

	var users []entities.User
	result := conn.Order("id desc").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func Create(u entities.User) (entities.User, error) {
	conn := db.Connection()
	_, err := ReadByEmail(u.Email)

	if err == nil {
		return entities.User{}, errors.New("email is already taken")
	}

	u.Password, _ = utils.ToHash(u.Password)

	result := conn.Create(&u)
	if result.Error != nil {
		return entities.User{}, errors.New("error creating user")
	}

	return u, nil
}

func Authenticate(u entities.User) (entities.User, error) {
	user, err := ReadByEmail(u.Email)

	if err != nil {
		return entities.User{}, errors.New("invalid credentials")
	}

	res := utils.Compare(user.Password, u.Password)

	if !res {
		return entities.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func ReadById(id uint) (entities.User, error) {
	return find("id = ?", id)
}

func ReadByEmail(email string) (entities.User, error) {
	return find("email = ?", email)
}

func SearchByEmail(email string) ([]entities.User, error) {
	//%%  a literal percent sign; consumes no value
	return findAll("email LIKE ?", fmt.Sprintf("%%%s%%", email))
}

func findAll(query interface{}, args interface{}) ([]entities.User, error) {
	conn := db.Connection()

	var users []entities.User

	result := conn.Where(query, args).Order("id desc").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func find(query interface{}, args interface{}) (entities.User, error) {
	conn := db.Connection()

	user := entities.User{}

	result := conn.Where(query, args).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func ReadTimeSeriesCount() (interface{}, error) {
	var users []struct {
		Date  int `json:"date"`
		Count int `json:"count"`
	}
	result := db.Connection().Raw(`
SELECT CONCAT(EXTRACT(YEAR_MONTH FROM created_at), EXTRACT(DAY FROM created_at)) as date, COUNT(id) as count
	FROM users
	GROUP BY date
	ORDER BY date DESC`).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func ReadSummaries() (interface{}, error) {
	var users struct {
		Today int `json:"today"`
		Month int `json:"month"`
		Year  int `json:"year"`
		Total int `json:"total"`
	}
	now := time.Now().UTC()
	today := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	month := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), 1)
	year := fmt.Sprintf("%d-%d-%d", now.Year(), 1, 1)

	result := db.Connection().Raw(`SELECT 
    	SUM(created_at > ?) as today,
    	SUM(created_at > ?) as month,
    	SUM(created_at > ?) as year,
       COUNT(created_at) as total
FROM users`, today, month, year).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
