package repositories

import (
	"accounts.sidooh/clients"
	Account "accounts.sidooh/models/account"
	Invite "accounts.sidooh/models/invite"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

func Create(a Account.Model) (Account.Model, error) {
	//	Get Invite if exists
	invite, err := Invite.UnexpiredByPhone(a.Phone)
	if err != nil {
		fmt.Println("Invite not found for", a.Phone)
	} else {
		a.InviterID = invite.InviterID
	}

	//	Create Account
	account, err := Account.Create(a)
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
	account, err := Account.ById(id)
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
	if util.Compare(account.Pin.String, pin) {
		return nil
	} else {
		return errors.New("pin is incorrect")
	}
}

func SetPin(id uint, pin string) error {
	//	Get Account
	account, err := Account.ById(id)
	if err != nil {
		return errors.New("account not found")
	}

	//	Set Pin
	hashedPin, err := util.ToHash(pin)
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
	account, err := Account.ById(id)
	if err != nil {
		return false
	}

	//	Check Pin exists
	return account.Pin.Valid
}

func UpdateProfile(id uint, name string) (User.Model, error) {
	//	Get Account
	account, err := Account.ByIdWithUser(id)
	if err != nil {
		return User.Model{}, errors.New("invalid credentials")
	}

	if account.User == nil {
		var user = User.Model{
			Name:     name,
			Username: account.Phone,
			IdNumber: account.Phone,
			Email:    account.Phone + "@sidooh.net",
			Status:   constants.ACTIVE,
		}

		user, err := User.CreateUser(user)
		if err != nil {
			return User.Model{}, err
		}

		account.Update("user_id", strconv.Itoa(int(user.ID)))

		return user, nil
	} else {
		account.User.Update("name", name)

		return *account.User, nil
	}
}

func GetAccounts(withUser bool) (interface{}, error) {
	if withUser {
		return Account.AllWithUser()
	} else {
		return Account.All()
	}
}

func GetAccountById(id uint, withUser bool, withInvite bool) (interface{}, error) {
	if withUser && withInvite {
		type AccountWithUserAndInviter struct {
			Account.ModelWithUser
			Inviter *Account.ModelWithUser `json:"inviter"`
		}

		account, err := Account.ByIdWithUser(id)
		if err != nil {
			return nil, err
		}

		inviter, err := Account.ByIdWithUser(account.InviterID)
		if err != nil {
			return &AccountWithUserAndInviter{
				ModelWithUser: *account,
			}, nil
		}

		return &AccountWithUserAndInviter{
			ModelWithUser: *account,
			Inviter:       inviter,
		}, err
	} else if withUser {
		return Account.ByIdWithUser(id)
	} else {
		return Account.ById(id)
	}
}

func GetAccountByPhone(phone string, withUser bool) (interface{}, error) {
	if withUser {
		return Account.ByPhoneWithUser(phone)
	} else {
		return Account.ByPhone(phone)
	}
}

func ResetPin(id uint) error {
	//	Get Account
	account, err := Account.ByIdWithUser(id)
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
	return Account.TimeSeriesCount(limit)
}

func GetAccountsSummary() (interface{}, error) {
	return Account.Summaries()
}
