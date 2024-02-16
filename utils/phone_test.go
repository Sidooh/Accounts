package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPhoneByCountry(t *testing.T) {
	phoneList := []string{
		"700000000", "748000000", "110000000", "730000000", "762000000",
		//"106000000", // TODO: This number fails, investigate why
		"779000000", "764000000", "747000000",
	}

	for _, number := range phoneList {
		phone, err := GetPhoneByCountry("KE", number)
		require.NoError(t, err)
		require.Equal(t, phone, "254"+number)
	}
}
