package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashing(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := ToHash(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isValidPassword := Compare(hashedPassword, password)

	require.True(t, isValidPassword)
}
