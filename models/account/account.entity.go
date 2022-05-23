package account

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/models/user"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Model struct {
	models.ModelID

	Phone     string         `json:"phone" gorm:"uniqueIndex; size:16"`
	Active    bool           `json:"active"`
	Pin       sql.NullString `json:"-"`
	TelcoID   int            `json:"-"`
	InviterID sql.NullInt32  `json:"-"`
	UserID    uint           `json:"-"`

	models.ModelTimeStamps
}

type ModelWithUser struct {
	Model

	User user.Model `json:"user"`
}

type InviteModel struct {
	Model

	Level int `json:"level"`
}

func (Model) TableName() string {
	return "accounts"
}
func (ModelWithUser) TableName() string {
	return "accounts"
}

func All() ([]Model, error) {
	var accounts []Model
	result := db.Connection().Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
}

// TODO: Check whether using pointers here saves memory
func Create(a Model) (Model, error) {
	_, err := ByPhone(a.Phone)
	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := db.Connection().Omit("UserID").Create(&a)
	if result.Error != nil {
		return Model{}, errors.New("error creating account")
	}

	return a, nil
}

func ById(id uint) (Model, error) {
	return find("id = ?", id)
}

func ByIdWithUser(id uint) (ModelWithUser, error) {
	account := ModelWithUser{}

	result := db.Connection().Joins("User").First(&account, id)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func ByPhone(phone string) (Model, error) {
	return find("phone = ?", phone)
}

func ByPhoneWithUser(phone string) (ModelWithUser, error) {
	account := ModelWithUser{}

	result := db.Connection().Where("accounts.phone = ?", phone).Joins("User").First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func SearchByPhone(phone string) ([]Model, error) {
	//%%  a literal percent sign; consumes no value
	return findAll("phone LIKE ?", fmt.Sprintf("%%%s%%", phone))
}

func SearchByIdOrPhone(search string) (Model, error) {
	account := Model{}

	result := db.Connection().Where("id = ? OR phone = ?", search, search).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil

}

func findAll(query interface{}, args interface{}) ([]Model, error) {
	var accounts []Model

	result := db.Connection().Where(query, args).Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
}

func find(query interface{}, args ...interface{}) (Model, error) {
	account := Model{}

	result := db.Connection().Where(query, args).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func (a *Model) Save() *gorm.DB {
	return db.Connection().Save(&a)
}

func (a *Model) Update(column string, value string) *gorm.DB {
	return db.Connection().Model(&a).Update(column, value)
}

// Referral/Invite Queries

// Ancestors 1. Get ancestors
func Ancestors(id uint, levelLimit uint) ([]InviteModel, error) {
	conn := db.Connection()

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
	conn := db.Connection()

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
