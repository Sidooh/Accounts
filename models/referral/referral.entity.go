package referral

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/util/constants"
	"database/sql"
	"errors"
	"time"
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

func All(db *db.DB) ([]Model, error) {
	var referrals []Model
	result := db.Conn.Find(&referrals)
	if result.Error != nil {
		return referrals, result.Error
	}

	return referrals, nil
}

func Create(db *db.DB, r Model) (Model, error) {
	if r.AccountID == 0 {
		return Model{}, errors.New("AccountId is required")
	}
	if r.Status == "" {
		r.Status = constants.PENDING
	}

	_, err := ByPhone(r.RefereePhone)
	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := db.Conn.Create(&r)
	if result.Error != nil {
		return Model{}, errors.New("error creating referral")
	}

	return r, nil
}

func ById(id uint) (Model, error) {
	return find("id = ?", id)
}

func ByPhone(phone string) (Model, error) {
	return find("referee_phone = ?", phone)
}

func UnexpiredByPhone(db *db.DB, phone string) (Model, error) {
	referral := Model{}

	result := db.Conn.
		Where("referee_phone", phone).
		Where("status", constants.PENDING).
		Where("created_at > ?", time.Now().Add(-48*time.Hour)).
		First(&referral)

	if result.Error != nil {
		return referral, result.Error
	}

	return referral, nil
}

func find(query interface{}, args ...interface{}) (Model, error) {
	conn := db.NewConnection().Conn

	referral := Model{}

	result := conn.Where(query, args).First(&referral)
	if result.Error != nil {
		return referral, result.Error
	}

	return referral, nil
}

func (r *Model) Save(db *db.DB) interface{} {
	return db.Conn.Save(&r)
}

func RemoveExpired(db *db.DB) error {
	var expired []Model

	db.Conn.
		Where("status", constants.PENDING).
		Where("created_at < ?", time.Now().Add(-48*time.Hour)).
		Delete(&expired)

	return nil
}
