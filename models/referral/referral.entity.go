package referral

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

	RefereePhone string        `json:"phone"`
	Status       string        `json:"active"`
	AccountID    uint          `json:"account_id"`
	RefereeID    sql.NullInt32 `json:"-"`
}

func (Model) TableName() string {
	return "referrals"
}

func All() ([]Model, error) {
	conn := db.NewConnection()

	var referrals []Model
	result := conn.Find(&referrals)
	if result.Error != nil {
		return referrals, result.Error
	}

	return referrals, nil
}

func Create(r Model, phone string) (Model, error) {
	conn := db.NewConnection()
	_, err := ByPhone(nil, r.RefereePhone)

	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := conn.Omit(clause.Associations).Create(&r)
	if result.Error != nil {
		return Model{}, errors.New("error creating referral")
	}

	return r, nil
}

func ByPhone(conn *gorm.DB, phone string) (Model, error) {
	if conn == nil {
		conn = db.NewConnection()
	}
	return find(conn, "referee_phone = ?", phone)
}

func find(conn *gorm.DB, query interface{}, args interface{}) (Model, error) {
	referral := Model{}

	result := conn.Where(query, args).First(&referral)
	if result.Error != nil {
		return referral, result.Error
	}

	return referral, nil
}

func (r *Model) Save(conn *gorm.DB) interface{} {
	return conn.Save(&r)
}

func removeExpired() {
	//TODO: Remove expired referrals
}
