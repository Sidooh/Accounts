package utils

import (
	"accounts.sidooh/utils/cache"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	viper.Set("APP_ENV", "TEST")

	cache.Init()

	os.Exit(m.Run())
}

func TestGenerateOTP(t *testing.T) {
	otp := GenerateOTP("714611696")
	require.NotEmpty(t, otp)
}

func TestCheckOTP(t *testing.T) {
	otp := GenerateOTP("714611696")

	valid := CheckOTP("714611696", 1)
	require.False(t, valid)

	valid = CheckOTP("714611696", otp)
	require.True(t, valid)
}
