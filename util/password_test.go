package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashing(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := ToHashV3(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isValidPassword := CompareV3(hashedPassword, password)

	require.True(t, isValidPassword)
}
