package otp

import (
	"accounts.sidooh/pkg/clients"
	"accounts.sidooh/utils"
	"errors"
	"fmt"
)

func GenerateOTP(phone string) error {
	notifyClient := clients.GetNotifyClient()

	otp := utils.GenerateOTP(phone)

	// Send OTP to phone number
	message := fmt.Sprintf("S-%v is your verification code.", otp)

	err := notifyClient.SendSMS("DEFAULT", phone, message)

	return err
}

func ValidateOTP(phone string, otp int) error {
	if utils.CheckOTP(phone, otp) {
		return nil
	}

	return errors.New("invalid otp")
}
