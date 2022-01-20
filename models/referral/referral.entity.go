package referral

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"database/sql"
	"errors"
	"gorm.io/gorm/clause"
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

func All() ([]Model, error) {
	conn := db.NewConnection()

	var referrals []Model
	result := conn.Find(&referrals)
	if result.Error != nil {
		return referrals, result.Error
	}

	return referrals, nil
}

func Create(r Model) (Model, error) {
	if r.AccountID == 0 {
		return Model{}, errors.New("AccountId is required")
	}
	if r.Status == "" {
		r.Status = models.PENDING
	}

	conn := db.NewConnection()
	_, err := ByPhone(r.RefereePhone)

	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := conn.Omit(clause.Associations).Create(&r)
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

func find(query interface{}, args interface{}) (Model, error) {
	conn := db.NewConnection()

	referral := Model{}

	result := conn.Where(query, args).First(&referral)
	if result.Error != nil {
		return referral, result.Error
	}

	return referral, nil
}

func (r *Model) Save() interface{} {
	conn := db.NewConnection()
	return conn.Save(&r)
}

func RemoveExpired() error {
	conn := db.NewConnection()

	var expired []Model

	conn.
		Where("status", models.PENDING).
		Where("created_at < ?", time.Now().Add(-48*time.Hour)).
		Delete(&expired)

	return nil
}
