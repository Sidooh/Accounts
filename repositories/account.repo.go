package repositories

import (
	Account "accounts.sidooh/models/account"
	Invite "accounts.sidooh/models/invite"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"errors"
	"fmt"
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

func HasPin(id uint) error {
	//	Get Account
	account, err := Account.ById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	//	Check Pin exists
	if account.Pin.Valid {
		return nil
	}

	return errors.New("invalid credentials")
}

func UpdateProfile(id uint, name string) (User.Model, error) {
	//	Get Account
	account, err := Account.ByIdWithUser(id)
	if err != nil {
		return User.Model{}, errors.New("invalid credentials")
	}

	if account.User.ID != 0 {
		account.User.Name = name

		account.User.Save()
	} else {
		account.User.Name = name
		account.User.Email = account.Phone + "@sidooh.net"

		user, err := User.CreateUser(account.User)
		if err != nil {
			return User.Model{}, err
		}

		account.User = user
		account.UserID = user.ID
		account.Save()
	}

	return account.User, nil

}
