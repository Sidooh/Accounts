package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

// RandomBool generates a random boolean
func RandomBool() bool {
	return rand.Intn(2) >= 1
}

// RandomIntBetween generates a random integer between min and max
func RandomIntBetween(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomInt generates a random integer between min and max
func RandomInt(length int) int64 {
	max := int64(math.Pow10(length)) - 1
	min := int64(math.Pow10(length - 1))

	return RandomIntBetween(min, max)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(lowerCharSet)

	for i := 0; i < n; i++ {
		c := lowerCharSet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomSecureString(n int) string {
	var sb strings.Builder
	k := len(allCharSet)

	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		c := allCharSet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a random name
func RandomName() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomPhone generates a random phone
func RandomPhone() string {
	prefix := "7"
	if RandomBool() {
		prefix = "1"
	}

	return fmt.Sprintf("%s%s", prefix, strconv.FormatInt(RandomInt(8), 10))
}
