package db

import (
	"accounts.sidooh/util"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
)

type DB struct {
	Conn *gorm.DB
}

func NewConnection() *DB {
	dsn := viper.GetString("DB_DSN")
	env := strings.ToUpper(viper.GetString("APP_ENV"))
	if env != "TEST" && dsn == "" {
		panic("database connection not set")
	}

	var db = new(gorm.DB)
	var err = errors.New("")

	if env == "TEST" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		file := util.GetLogFile("db.log")

		newLogger := logger.New(
			log.New(file, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel: logger.Info, // Log level
			},
		)

		config := &gorm.Config{
			Logger: newLogger,
		}

		db, err = gorm.Open(mysql.Open(dsn), config)
	}

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	return &DB{Conn: db}
}
