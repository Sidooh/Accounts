package account

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	models.Model

	Phone      string         `json:"phone"`
	Active     bool           `json:"active"`
	Pin        sql.NullString `json:"-"`
	TelcoID    int            `json:"-"`
	ReferrerID sql.NullInt32  `json:"-"`
	UserID     sql.NullInt32  `json:"-"`
}

func (Model) TableName() string {
	return "accounts"
}

func All(db *db.DB) ([]Model, error) {
	var accounts []Model
	result := db.Conn.Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
}

func Create(db *db.DB, a Model) (Model, error) {
	_, err := ByPhone(db, a.Phone)
	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := db.Conn.Omit(clause.Associations).Create(&a)
	if result.Error != nil {
		fmt.Println(result.Error)
		return Model{}, errors.New("error creating account")
	}

	return a, nil
}

func ById(db *db.DB, id uint) (Model, error) {
	return find(db, "id = ?", id)
}

func ByPhone(db *db.DB, phone string) (Model, error) {
	return find(db, "phone = ?", phone)
}

func find(db *db.DB, query interface{}, args interface{}) (Model, error) {
	account := Model{}

	result := db.Conn.Where(query, args).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func (a *Model) Save(db *db.DB) *gorm.DB {
	return db.Conn.Save(&a)
}
