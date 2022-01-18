package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomBool generates a random boolean
func RandomBool() bool {
	return rand.Intn(2) >= 1
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
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

	return fmt.Sprintf("%s%s", prefix, strconv.FormatInt(RandomInt(10000000, 99999999), 10))
}
