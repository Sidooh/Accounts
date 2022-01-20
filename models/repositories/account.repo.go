package repositories

import (
	Account "accounts.sidooh/models/account"
	Referral "accounts.sidooh/models/referral"
	"database/sql"
	"fmt"
)

func Create(a Account.Model) (Account.Model, error) {
	//	Get Referral if exists
	referral, err := Referral.ByPhone(a.Phone)
	if err != nil {
		fmt.Println("Referral not found for", a.Phone)
	} else {
		a.ReferrerID = sql.NullInt32{
			Int32: int32(referral.AccountID),
			Valid: true,
		}
	}

	//	Create Account
	account, err := Account.Create(a)
	if err != nil {
		return a, err
	}

	//	Update referral
	if referral.ID != 0 {

		referral.RefereeID = sql.NullInt32{
			Int32: int32(account.ID),
			Valid: true,
		}
		referral.Status = "active"
		referral.Save()
	}

	return account, nil
}
