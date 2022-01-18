package account

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	models.Model

	Phone      string         `json:"phone"`
	Active     bool           `json:"active"`
	Pin        sql.NullString `json:"-"`
	TelcoID    int            `json:"telco_id"`
	ReferrerID sql.NullInt32  `json:"-"`
	UserID     sql.NullInt32  `json:"-"`
}

func (Model) TableName() string {
	return "accounts"
}

func All() ([]Model, error) {
	conn := db.NewConnection()

	var accounts []Model
	result := conn.Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
}

func Create(conn *gorm.DB, a Model) (Model, error) {
	if conn == nil {
		conn = db.NewConnection()
	}
	_, err := ByPhone(a.Phone)
	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := conn.Omit(clause.Associations).Create(&a)
	if result.Error != nil {
		return Model{}, errors.New("error creating account")
	}

	return a, nil
}

func ById(id uint) (Model, error) {
	return find("id = ?", id)
}

func ByPhone(phone string) (Model, error) {
	return find("phone = ?", phone)
}

func find(query interface{}, args interface{}) (Model, error) {
	conn := db.NewConnection()

	account := Model{}

	result := conn.Where(query, args).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}
