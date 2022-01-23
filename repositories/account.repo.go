package repositories

import (
	"accounts.sidooh/db"
	Account "accounts.sidooh/models/account"
	Referral "accounts.sidooh/models/referral"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"database/sql"
	"errors"
	"fmt"
)

var datastore = new(db.DB)

func Construct(db *db.DB) {
	datastore = db
}

func Create(a Account.Model) (Account.Model, error) {
	//	Get Referral if exists
	referral, err := Referral.UnexpiredByPhone(datastore, a.Phone)
	if err != nil {
		fmt.Println("Referral not found for", a.Phone)
	} else {
		a.ReferrerID = sql.NullInt32{
			Int32: int32(referral.AccountID),
			Valid: true,
		}
	}

	//	Create Account
	account, err := Account.Create(datastore, a)
	if err != nil {
		return a, err
	}

	//	Update referral
	if referral.ID != 0 {

		referral.RefereeID = sql.NullInt32{
			Int32: int32(account.ID),
			Valid: true,
		}
		referral.Status = constants.ACTIVE
		referral.Save(datastore)
	}

	return account, nil
}

func CheckPin(id uint, pin string) error {
	//	Get Account
	account, err := Account.ById(datastore, id)
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
	account, err := Account.ById(datastore, id)
	if err != nil {
		return errors.New("account not found")
	}

	//	Set Pin
	hashedPin, err := util.ToHash(pin)
	if err != nil {
		return errors.New("unable to set pin")
	}

	result := account.Update(datastore, "pin", hashedPin)
	if result.Error != nil {
		return errors.New("unable to set pin")
	}

	return nil
}
