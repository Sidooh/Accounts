package util

import (
	"errors"
	"github.com/ttacon/libphonenumber"
	"regexp"
	"strings"
)

func GetPhoneByCountry(country string, phone string) (string, error) {
	// TODO: Check out and substitute https://pkg.go.dev/github.com/nyaruka/phonenumbers

	num, err := libphonenumber.Parse(phone, country)
	if err != nil {
		return phone, err
	}

	valid := libphonenumber.IsValidNumber(num) || phoneValidKE(phone)
	if !valid {
		return phone, errors.New("number is not valid")
	}

	phone = strings.TrimPrefix(libphonenumber.Format(num, libphonenumber.E164), "+")

	return phone, nil
}

func phoneValidKE(phone string) bool {
	matchString, err := regexp.MatchString("/^(?:254|\\+254|0)?((?:7(?:[0129][0-9]|4[0123568]|5[789]|6[89])|(1([1][0-5])))[0-9]{6})$/", phone)
	if err != nil {
		return false
	}

	return matchString
}
