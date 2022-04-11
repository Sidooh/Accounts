package account

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/models/user"
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

type ModelWithUser struct {
	Model

	User user.User `json:"user"`
}

type InviteModel struct {
	Model

	Level int `json:"level"`
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

func ByIdWithUser(db *db.DB, id uint) (ModelWithUser, error) {
	account := ModelWithUser{}

	result := db.Conn.Joins("User").First(&account, id)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func ByPhone(db *db.DB, phone string) (Model, error) {
	return find(db, "phone = ?", phone)
}

func ByPhoneWithUser(db *db.DB, phone string) (ModelWithUser, error) {
	account := ModelWithUser{}

	result := db.Conn.Where("accounts.phone = ?", phone).Joins("User").First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func SearchByPhone(db *db.DB, phone string) ([]Model, error) {
	//%%  a literal percent sign; consumes no value
	return findAll(db, "phone LIKE ?", fmt.Sprintf("%%%s%%", phone))
}

func findAll(db *db.DB, query interface{}, args interface{}) ([]Model, error) {
	var accounts []Model

	result := db.Conn.Where(query, args).Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
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

func (a *Model) Update(db *db.DB, column string, value string) *gorm.DB {
	return db.Conn.Model(&a).Update(column, value)
}

// Referral/Invite Queries

// Ancestors 1. Get ancestors
func Ancestors(id uint, levelLimit uint) ([]InviteModel, error) {
	conn := db.NewConnection().Conn

	var accounts []InviteModel
	conn.Raw(
		"WITH RECURSIVE ancestors (id, phone, referrer_id, level) AS\n"+
			"("+
			"SELECT id, phone, referrer_id, 0 level\n"+
			"FROM accounts\n"+
			"WHERE id = ?\n"+
			"UNION ALL\n"+
			"SELECT a.id, a.phone, a.referrer_id, an.level+1\n"+
			"FROM ancestors AS an JOIN accounts AS a\n"+
			"ON an.referrer_id = a.id"+
			")\n"+
			"SELECT * FROM ancestors LIMIT ?",
		id, levelLimit).
		Scan(&accounts)

	return accounts, nil
}

// Descendants 2. Get descendants
func Descendants(id uint, levelLimit uint) ([]InviteModel, error) {
	conn := db.NewConnection().Conn

	var accounts []InviteModel

	conn.Raw(
		"WITH RECURSIVE descendants (id, phone, referrer_id, level) AS\n"+
			"("+
			"SELECT id, phone, referrer_id, 0 level\n"+
			"FROM accounts\n"+
			"WHERE id = ?\n"+
			"UNION ALL\n"+
			"SELECT a.id, a.phone, a.referrer_id, d.level+1\n"+
			"FROM descendants AS d JOIN accounts AS a\n"+
			"ON d.id = a.referrer_id"+
			")\n"+
			"SELECT * FROM descendants WHERE level < ?",
		id, levelLimit).
		Scan(&accounts)

	return accounts, nil
}
