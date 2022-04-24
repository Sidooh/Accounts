package invite

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/models/account"
	"accounts.sidooh/util/constants"
	"errors"
	"github.com/spf13/viper"
	"time"
)

type Model struct {
	models.ModelID

	Phone     string `json:"phone" gorm:"uniqueIndex; size:16"`
	Status    string `json:"status" gorm:"size:16"`
	AccountID uint   `json:"account_id,omitempty"`
	InviterID uint   `json:"inviter_id"`

	models.ModelTimeStamps
}

type ModelWithAccountAndInvite struct {
	Model

	//TODO: Add a constraint to ensure these 2 can't have same values
	// 	i.e. a user can't invite themselves, obviously
	Account account.Model `json:"account"`
	Inviter account.Model `json:"inviter"`
}

func (Model) TableName() string {
	return "invites"
}

func (ModelWithAccountAndInvite) TableName() string {
	return "invites"
}

func All() ([]Model, error) {
	var invites []Model
	result := db.Connection().Find(&invites)
	if result.Error != nil {
		return invites, result.Error
	}

	return invites, nil
}

func Create(r Model) (Model, error) {
	if r.InviterID == 0 {
		return Model{}, errors.New("inviter_id is required")
	}
	if r.Status == "" {
		r.Status = constants.PENDING
	}

	_, err := ByPhone(r.Phone)
	if err == nil {
		return Model{}, errors.New("phone is already taken")
	}

	result := db.Connection().Omit("AccountID").Create(&r)
	if result.Error != nil {
		return Model{}, errors.New("error creating invite")
	}

	return r, nil
}

func ById(id uint) (Model, error) {
	return find("id = ?", id)
}

func ByPhone(phone string) (Model, error) {
	return find("phone = ?", phone)
}

func UnexpiredByPhone(phone string) (Model, error) {
	//TODO: Move the defaults to Config struct and remove from file
	expiryTime := time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

	invite := Model{}

	result := db.Connection().
		Where("phone", phone).
		Where("status", constants.PENDING).
		Where("created_at > ?", time.Now().Add(-expiryTime)).
		First(&invite)

	if result.Error != nil {
		return invite, result.Error
	}

	return invite, nil
}

func Unexpired() ([]Model, error) {
	var invites []Model

	result := db.Connection().
		Where("status <> ?", constants.EXPIRED).
		Find(&invites)

	if result.Error != nil {
		return invites, result.Error
	}

	return invites, nil
}

func find(query interface{}, args ...interface{}) (Model, error) {
	conn := db.Connection()

	invite := Model{}

	result := conn.Where(query, args).First(&invite)
	if result.Error != nil {
		return invite, result.Error
	}

	return invite, nil
}

func (r *Model) Save() interface{} {
	return db.Connection().Save(&r)
}

func MarkExpired() error {
	expiryTime := time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

	db.Connection().
		Model(&Model{}).
		Where("status", constants.PENDING).
		Where("created_at < ?", time.Now().Add(-expiryTime)).
		Update("status", constants.EXPIRED)

	return nil
}
