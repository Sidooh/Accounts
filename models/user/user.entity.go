package user

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/models/account"
	"accounts.sidooh/services"
	"errors"
	"github.com/SamuelTissot/sqltime"
)

type User struct {
	models.Model

	Email           string       `json:"email"`
	Password        string       `json:"-"`
	Name            string       `json:"name"`
	Username        string       `json:"username"`
	IdNumber        string       `json:"id_number"`
	Status          string       `json:"status"`
	EmailVerifiedAt sqltime.Time `json:"-"`

	Account account.Model
}

func CreateUser(u User) (User, error) {
	conn := db.NewConnection()
	user, err := findUserByEmail(u.Email)

	if err == nil {
		return User{}, errors.New("email is already taken")
	}

	u.Password, _ = services.ToHash(u.Password)

	result := conn.Create(&u)
	if result.Error != nil {
		return User{}, errors.New("error creating user")
	}

	//user, _ = findUserById(u.ID)

	return user, nil
}

func AuthUser(u User) (User, error) {
	user, err := findUserByEmail(u.Email)

	if err != nil {
		return User{}, errors.New("invalid credentials")
	}

	res := services.Compare(user.Password, u.Password)

	if !res {
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func findUserById(id uint) (User, error) {
	return findUser("id = ?", id)
}

func findUserByEmail(email string) (User, error) {
	return findUser("email = ?", email)
}

func findUser(query interface{}, args interface{}) (User, error) {
	conn := db.NewConnection()

	user := User{}

	result := conn.Where(query, args).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
