package accounts

import (
	"accounts.sidooh/pkg/clients"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/pkg/repositories/invites"
	"accounts.sidooh/pkg/repositories/users"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

func Create(a entities.Account) (entities.Account, error) {
	//	Get Invite if exists
	invite, err := invites.ReadUnexpiredByPhone(a.Phone)
	if err != nil {
		fmt.Println("Invite not found for", a.Phone)
	} else {
		a.InviterID = invite.InviterID
	}

	//	Create Account
	account, err := CreateAccount(a)
	if err != nil {
		return a, err
	}

	//	Update invite
	if invite.ID != 0 {

		invite.AccountID = account.ID
		invite.Status = constants.ACTIVE
		invite.Save()
	}

	return account, nil
}

func CheckPin(id uint, pin string) error {
	//	Get Account
	account, err := ReadById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	//	Check Pin

	// Account for existing older pins
	if len(account.Pin.String) == 4 {
		if account.Pin.String == pin {
			err := SetPin(account.ID, account.Pin.String)
			if err != nil {
				//Log error
				return err
			}
			return nil
		} else {
			return errors.New("invalid credentials")
		}
	}

	// New algorithm
	if utils.Compare(account.Pin.String, pin) {
		return nil
	} else {
		return errors.New("pin is incorrect")
	}
}

func SetPin(id uint, pin string) error {
	//	Get Account
	account, err := ReadById(id)
	if err != nil {
		return errors.New("account not found")
	}

	//	Set Pin
	hashedPin, err := utils.ToHash(pin)
	if err != nil {
		return errors.New("unable to set pin")
	}

	result := account.Update("pin", hashedPin)
	if result.Error != nil {
		return errors.New("unable to set pin")
	}

	return nil
}

func HasPin(id uint) bool {
	//	Get Account
	account, err := ReadById(id)
	if err != nil {
		return false
	}

	//	Check Pin exists
	return account.Pin.Valid
}

func UpdateProfile(id uint, name string) (entities.User, error) {
	//	Get Account
	account, err := ReadWithUser(id)
	if err != nil {
		return entities.User{}, errors.New("invalid credentials")
	}

	if account.User == nil {
		var user = entities.User{
			Name:     name,
			Username: account.Phone,
			IdNumber: account.Phone,
			Email:    account.Phone + "@sidooh.net",
			Status:   constants.ACTIVE,
		}

		user, err := users.Create(user)
		if err != nil {
			return entities.User{}, err
		}

		account.Update("user_id", strconv.Itoa(int(user.ID)))

		return user, nil
	} else {
		account.User.Update("name", name)

		return *account.User, nil
	}
}

func GetAccounts(withUser bool, withInviter bool, limit int) (interface{}, error) {
	if withUser && withInviter {
		return ReadAllWithUserAndInviter(limit)
	} else if withUser {
		return ReadAllWithUser(limit)
	} else {
		return ReadAll()
	}
}

func GetAccountById(id uint, withUser bool, withInviter bool) (interface{}, error) {
	if withUser && withInviter {
		return ReadWithUserAndInviter(id)
	} else if withUser {
		return ReadWithUser(id)
	} else {
		return ReadById(id)
	}
}

func GetAccountByPhone(phone string, withUser bool) (interface{}, error) {
	if withUser {
		return ReadByPhoneWithUser(phone)
	} else {
		return ReadByPhone(phone)
	}
}

func ResetPin(id uint) error {
	//	Get Account
	account, err := ReadWithUser(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	account.Pin = sql.NullString{}
	account.Save()

	notifyClient := clients.GetNotifyClient()

	name := ""
	if account.User != nil {
		name = " " + account.User.Name
	}

	message := fmt.Sprintf("Hello%v,\nYour Sidooh pin has been reset.\n"+
		"Dial *384*99# to set a new pin and keep earning from your purchases.\n\n"+
		"If you did not request a pin reset, kindly contact us at customersupport@sidooh.co.ke promptly!", name)
	if err := notifyClient.SendSMS("DEFAULT", account.Phone, message); err != nil {
		return err
	}

	return nil
}

func GetAccountsTimeData(limit int) (interface{}, error) {
	return ReadAccountsTimeSeriesCount(limit)
}

func GetAccountsSummary() (interface{}, error) {
	return ReadAccountsSummaries()
}
