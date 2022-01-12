package services

import (
	"golang.org/x/crypto/bcrypt"
)

func ToHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)

	return string(bytes), err
}

func Compare(storedPassword string, suppliedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(suppliedPassword))

	return err == nil
}

//func ToHash(password string) (string, error) {
//	salt := make([]byte, 8)
//	_, err := rand.Read(salt)
//	if err != nil {
//		return "", err
//	}
//
//	buf, err := scrypt.Key([]byte(password), salt, 32768, 12, 1, 64)
//	if err != nil {
//		return "", err
//	}
//
//	return fmt.Sprintf("%v.%v", hex.EncodeToString(buf), hex.EncodeToString(salt)), nil
//}
//
//
//func Compare(storedPassword string, suppliedPassword string) bool {
//	split := strings.Split(storedPassword, ".")
//
//	salt, _ := hex.DecodeString(split[1])
//	buf, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 12, 1, 64)
//	if err != nil {
//		return false
//	}
//
//	return hex.EncodeToString(buf) == split[0]
//}
