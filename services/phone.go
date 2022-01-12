package services

import (
	"errors"
	"fmt"
	"github.com/ttacon/libphonenumber"
	"strings"
)

func GetPhoneByCountry(country string, phone string) (string, error) {
	num, err := libphonenumber.Parse(phone, country)
	if err != nil {
		return phone, err
	}

	valid := libphonenumber.IsValidNumber(num)
	if !valid {
		return phone, errors.New("number is not valid")
	}

	fmt.Println(num)
	phone = strings.TrimPrefix(libphonenumber.Format(num, libphonenumber.E164), "+")

	return phone, nil
}
