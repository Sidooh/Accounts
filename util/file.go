package util

import (
	"fmt"
	"os"
)

func GetFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return file
}
