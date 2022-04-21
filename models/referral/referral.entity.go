package referral

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/util/constants"
	"database/sql"
	"errors"
	"github.com/spf13/viper"
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

//TODO: Move the defaults to Config struct and remove from file
var expiryTime = time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

func All() ([]Model, error) {
	var referrals []Model
	result := db.Connection().Find(&referrals)
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
		r.Status = constants.PENDING
	}

	_, err := ByPhone(r.RefereePhone)
	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := db.Connection().Create(&r)
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

func UnexpiredByPhone(phone string) (Model, error) {
	//TODO: Move the defaults to Config struct and remove from file
	expiryTime = time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

	referral := Model{}

	result := db.Connection().
		Where("referee_phone", phone).
		Where("status", constants.PENDING).
		Where("created_at > ?", time.Now().Add(-expiryTime)).
		First(&referral)

	if result.Error != nil {
		return referral, result.Error
	}

	return referral, nil
}

func Unexpired() ([]Model, error) {
	var referrals []Model

	result := db.Connection().
		Where("status <> ?", constants.EXPIRED).
		Find(&referrals)

	if result.Error != nil {
		return referrals, result.Error
	}

	return referrals, nil
}

func find(query interface{}, args ...interface{}) (Model, error) {
	conn := db.Connection()

	referral := Model{}

	result := conn.Where(query, args).First(&referral)
	if result.Error != nil {
		return referral, result.Error
	}

	return referral, nil
}

func (r *Model) Save() interface{} {
	return db.Connection().Save(&r)
}

func MarkExpired() error {
	expiryTime = time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

	db.Connection().
		Model(&Model{}).
		Where("status", constants.PENDING).
		Where("created_at < ?", time.Now().Add(-expiryTime)).
		Update("status", constants.EXPIRED)

	return nil
}
