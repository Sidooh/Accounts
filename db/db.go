package db

import (
	"errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection() *gorm.DB {
	dsn := viper.GetString("DB_DSN")

	env := viper.GetString("APP_ENV")

	var db = new(gorm.DB)
	var err = errors.New("")

	if env == "TEST" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		config := &gorm.Config{}

		if env == "PRODUCTION" {
			// TODO: Add db logger configs to redirect sql logs to file
			//		config = &gorm.Config{}
		}

		db, err = gorm.Open(mysql.Open(dsn), config)
	}

	if err != nil {
		panic(err)
	}

	return db
}
