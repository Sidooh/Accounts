package user

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/models/account"
	"accounts.sidooh/util"
	"database/sql"
	"errors"
)

type User struct {
	models.Model

	Email           string       `json:"email"`
	Password        string       `json:"-"`
	Name            string       `json:"name"`
	Username        string       `json:"username"`
	IdNumber        string       `json:"id_number"`
	Status          string       `json:"status"`
	EmailVerifiedAt sql.NullTime `json:"-"`

	Account account.Model `json:"-"`
}

func CreateUser(u User) (User, error) {
	conn := db.NewConnection().Conn
	_, err := FindUserByEmail(u.Email)

	if err == nil {
		return User{}, errors.New("email is already taken")
	}

	u.Password, _ = util.ToHash(u.Password)

	result := conn.Create(&u)
	if result.Error != nil {
		return User{}, errors.New("error creating user")
	}

	return u, nil
}

func AuthUser(u User) (User, error) {
	user, err := FindUserByEmail(u.Email)

	if err != nil {
		return User{}, errors.New("invalid credentials")
	}

	res := util.Compare(user.Password, u.Password)

	if !res {
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	return find("id = ?", id)
}

func FindUserByEmail(email string) (User, error) {
	return find("email = ?", email)
}

func find(query interface{}, args interface{}) (User, error) {
	conn := db.NewConnection().Conn

	user := User{}

	result := conn.Where(query, args).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
