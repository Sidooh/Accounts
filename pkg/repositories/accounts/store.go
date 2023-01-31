package accounts

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/utils"
	"errors"
	"fmt"
	"time"
)

func ReadAll() ([]entities.Account, error) {
	var accounts []entities.Account
	result := db.Connection().Order("id desc").Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
}

func ReadAllWithUser(limit int) ([]interface{}, error) {
	var accountsWithUsers []entities.AccountWithUser
	query := db.Connection().Joins("User").Order("id desc")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&accountsWithUsers)

	if result.Error != nil {
		return nil, result.Error
	}

	var accounts []interface{}
	for _, accountWithUser := range accountsWithUsers {
		if accountWithUser.UserID == 0 {
			accountModel := new(entities.Account)
			utils.ConvertStruct(accountWithUser, accountModel)
			accounts = append(accounts, accountModel)
		} else {
			accounts = append(accounts, accountWithUser)
		}
	}

	return accounts, nil
}

func ReadAllWithUserAndInviter(limit int) ([]entities.AccountWithUserAndInviter, error) {
	var accountsWithUserAndInviters []entities.AccountWithUserAndInviter
	query := db.Connection().Preload("User").Preload("Inviter").Preload("Inviter.User").Order("id desc")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&accountsWithUserAndInviters)

	if result.Error != nil {
		return nil, result.Error
	}

	return accountsWithUserAndInviters, nil
}

// TODO: Check whether using pointers here saves memory
func CreateAccount(a entities.Account) (entities.Account, error) {
	_, err := ReadByPhone(a.Phone)
	if err == nil {
		return entities.Account{}, errors.New("phone is already taken")
	}

	result := db.Connection().Omit("UserID").Create(&a)
	if result.Error != nil {
		return entities.Account{}, errors.New("error creating account")
	}

	return a, nil
}

func ReadById(id uint) (entities.Account, error) {
	return find("id = ?", id)
}

func ReadWithUser(id uint) (*entities.AccountWithUser, error) {
	accountWithUser := new(entities.AccountWithUser)

	result := db.Connection().Joins("User").First(&accountWithUser, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return accountWithUser, nil
}

func ReadWithUserAndInviter(id uint) (*entities.AccountWithUserAndInviter, error) {
	accountWithUserAndInviter := new(entities.AccountWithUserAndInviter)

	result := db.Connection().Joins("User").Joins("Inviter").First(&accountWithUserAndInviter, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return accountWithUserAndInviter, nil
}

func ReadByPhone(phone string) (entities.Account, error) {
	return find("phone = ?", phone)
}

func ReadByPhoneWithUser(phone string) (entities.AccountWithUser, error) {
	account := entities.AccountWithUser{}

	result := db.Connection().Where("accounts.phone = ?", phone).Joins("User").First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

func SearchByPhone(phone string) ([]entities.Account, error) {
	//%%  a literal percent sign; consumes no value
	return findMany("phone LIKE ?", fmt.Sprintf("%%%s%%", phone))
}

func SearchByIdOrPhone(search string) (entities.Account, error) {
	account := entities.Account{}

	result := db.Connection().Where("id = ? OR phone = ?", search, search).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil

}

func findMany(query interface{}, args interface{}) ([]entities.Account, error) {
	var accounts []entities.Account

	result := db.Connection().Where(query, args).Order("id desc").Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}

	return accounts, nil
}

func find(query interface{}, args ...interface{}) (entities.Account, error) {
	account := entities.Account{}

	result := db.Connection().Where(query, args).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

// Referral/Invite Queries

// Ancestors 1. Get ancestors
func ReadAncestors(id uint, levelLimit uint) ([]entities.InviteModel, error) {
	conn := db.Connection()

	var accounts []entities.InviteModel
	conn.Raw(
		"WITH RECURSIVE ancestors (id, phone, inviter_id, level) AS\n"+
			"("+
			"SELECT id, phone, inviter_id, 0 level\n"+
			"FROM accounts\n"+
			"WHERE id = ?\n"+
			"UNION ALL\n"+
			"SELECT a.id, a.phone, a.inviter_id, an.level+1\n"+
			"FROM ancestors AS an JOIN accounts AS a\n"+
			"ON an.inviter_id = a.id"+
			")\n"+
			"SELECT * FROM ancestors LIMIT ?",
		id, levelLimit).
		Scan(&accounts)

	if len(accounts) == 0 {
		return nil, errors.New("record not found")
	}

	return accounts, nil
}

// Descendants 2. Get descendants
func ReadDescendants(id uint, levelLimit uint) ([]entities.InviteModel, error) {
	conn := db.Connection()

	var accounts []entities.InviteModel

	conn.Raw(
		"WITH RECURSIVE descendants (id, phone, inviter_id, level) AS\n"+
			"("+
			"SELECT id, phone, inviter_id, 0 level\n"+
			"FROM accounts\n"+
			"WHERE id = ?\n"+
			"UNION ALL\n"+
			"SELECT a.id, a.phone, a.inviter_id, d.level+1\n"+
			"FROM descendants AS d JOIN accounts AS a\n"+
			"ON d.id = a.inviter_id"+
			")\n"+
			"SELECT * FROM descendants WHERE level < ?",
		id, levelLimit).
		Scan(&accounts)

	if len(accounts) == 0 {
		return nil, errors.New("record not found")
	}

	return accounts, nil
}

func ReadAccountsTimeSeriesCount(limit int) (interface{}, error) {
	var accounts []struct {
		Date  int `json:"date"`
		Count int `json:"count"`
	}
	result := db.Connection().Raw(`
SELECT EXTRACT(YEAR_MONTH FROM created_at) as date, COUNT(id) as count
	FROM accounts
	GROUP BY date
	ORDER BY date DESC
	LIMIT ?`, limit).Scan(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func ReadAccountsSummaries() (interface{}, error) {
	var accounts struct {
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
FROM accounts`, today, month, year).Scan(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}
