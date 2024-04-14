package utils

import (
	"accounts.sidooh/utils/cache"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func GenerateOTP(phone string) int {

	otp := RandomInt(6)

	SetOTP(phone, int(otp))

	return int(otp)
}

func CheckOTP(key string, otp int) bool {
	savedOtp := cache.Cache.Get(fmt.Sprintf("otp_%s", key))
	return savedOtp != nil && (*savedOtp).(int) == otp
}

func SetOTP(key string, otp int) {
	duration := viper.GetInt("OTP_VALIDITY")
	if duration == 0 {
		duration = 120
	}

	cache.Cache.Set(fmt.Sprintf("otp_%s", key), otp, time.Duration(duration)*time.Second)
}
