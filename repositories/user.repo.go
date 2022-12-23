package repositories

import (
	"accounts.sidooh/clients"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util"
	"errors"
	"fmt"
)

func ResetPassword(id uint) error {
	user, err := User.ById(id)
	if err != nil {
		return errors.New("invalid credentials")
	}

	newPassword := util.RandomSecureString(10)

	user.Password, _ = util.ToHash(newPassword)
	user.Save()

	notifyClient := clients.GetNotifyClient()

	message := fmt.Sprintf("Hello %v,\n Your account password has been reset. Below is your new default password: \n\n "+
		"%v \n\n Please ensure to change it after your next login.", user.Name, newPassword)
	if err := notifyClient.SendMail("DEFAULT", user.Email, message); err != nil {
		return err
	}

	return nil
}
