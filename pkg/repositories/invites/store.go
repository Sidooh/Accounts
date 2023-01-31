package invites

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/utils/constants"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func ReadAll(limit int) ([]entities.Invite, error) {
	var invites []entities.Invite
	query := db.Connection().Order("id desc")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&invites)
	if result.Error != nil {
		return invites, result.Error
	}

	return invites, nil
}
func ReadAllWithInviter(limit int) ([]entities.InviteWithInviter, error) {
	var invites []entities.InviteWithInviter
	query := db.Connection().Preload("Inviter").Preload("Inviter.User").Order("id desc")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&invites)
	if result.Error != nil {
		return invites, result.Error
	}

	return invites, nil
}

func Create(r entities.Invite) (entities.Invite, error) {
	if r.InviterID == 0 {
		return entities.Invite{}, errors.New("inviter_id is required")
	}
	if r.Status == "" {
		r.Status = constants.PENDING
	}

	_, err := ReadByPhone(r.Phone)
	if err == nil {
		return entities.Invite{}, errors.New("phone is already taken")
	}

	result := db.Connection().Omit("AccountID").Create(&r)
	if result.Error != nil {
		fmt.Println(result.Error)
		return entities.Invite{}, errors.New("error creating invite")
	}

	return r, nil
}

func ReadById(id uint) (entities.Invite, error) {
	return find("id = ?", id)
}

func ReadWithAccount(id uint) (*entities.InviteWithAccount, error) {
	inviteWithAccount := new(entities.InviteWithAccount)

	result := db.Connection().Joins("Account").First(&inviteWithAccount, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return inviteWithAccount, nil
}

func ReadWithAccountAndInviter(id uint) (*entities.InviteWithAccountAndInviter, error) {
	inviteWithAccountAndInviter := new(entities.InviteWithAccountAndInviter)

	result := db.Connection().Joins("Account").Joins("Inviter").First(&inviteWithAccountAndInviter, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return inviteWithAccountAndInviter, nil
}

func ReadByPhone(phone string) (entities.Invite, error) {
	return find("phone = ?", phone)
}

func ReadUnexpiredByPhone(phone string) (entities.Invite, error) {
	//TODO: Move the defaults to Config struct and remove from file
	expiryTime := time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

	invite := entities.Invite{}

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

func ReadUnexpired() ([]entities.Invite, error) {
	var invites []entities.Invite

	result := db.Connection().
		Where("status <> ?", constants.EXPIRED).
		Find(&invites)

	if result.Error != nil {
		return invites, result.Error
	}

	return invites, nil
}

func find(query interface{}, args ...interface{}) (entities.Invite, error) {
	conn := db.Connection()

	invite := entities.Invite{}

	result := conn.Where(query, args).First(&invite)
	if result.Error != nil {
		return invite, result.Error
	}

	return invite, nil
}

func MarkExpired() error {
	expiryTime := time.Duration(viper.GetFloat64("INVITE_EXPIRY")) * time.Hour

	db.Connection().
		Model(&entities.Invite{}).
		Where("status", constants.PENDING).
		Where("created_at < ?", time.Now().Add(-expiryTime)).
		Update("status", constants.EXPIRED)

	return nil
}

func ReadTimeSeriesCount(limit int) (interface{}, error) {
	var invites []struct {
		Date  int `json:"date"`
		Count int `json:"count"`
	}
	result := db.Connection().Raw(`
SELECT EXTRACT(YEAR_MONTH FROM created_at) as date, COUNT(id) as count
	FROM invites
	GROUP BY date
	ORDER BY date DESC
	LIMIT ?`, limit).Scan(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}

func ReadSummaries() (interface{}, error) {
	var invites struct {
		Today int `json:"today"`
		Month int `json:"month"`
		Year  int `json:"year"`
		Total int `json:"total"`
	}
	now := time.Now().UTC()
	today := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	month := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), 1)
	year := fmt.Sprintf("%d-%d-%d", now.Year(), 1, 1)

	result := db.Connection().Raw(`SELECT 
    	SUM(created_at > ?) as today,
    	SUM(created_at > ?) as month,
    	SUM(created_at > ?) as year,
       COUNT(created_at) as total
FROM invites`, today, month, year).Scan(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}
