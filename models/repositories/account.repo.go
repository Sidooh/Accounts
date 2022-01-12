package repositories

import (
	"accounts.sidooh/db"
	Account "accounts.sidooh/models/account"
	Referral "accounts.sidooh/models/referral"
	"database/sql"
	"fmt"
)

func Create(a Account.Model) (Account.Model, error) {
	conn := db.NewConnection()

	//	Get Referral if exists
	referral, err := Referral.ByPhone(conn, a.Phone)
	if err != nil {
		fmt.Println("Referral not found.")
	}

	//	Create Account
	a.ReferrerID = sql.NullInt32{
		Int32: int32(referral.AccountID),
		Valid: true,
	}
	account, err := Account.Create(conn, a)
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
		referral.Save(conn)
	}

	return account, nil
}
