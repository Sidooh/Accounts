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

var db *DB

func Connection() *gorm.DB {
	return db.Conn
}

func Init() {
	fmt.Println("==== Initializing DB ====")

	dsn := viper.GetString("DB_DSN")
	env := strings.ToUpper(viper.GetString("APP_ENV"))
	if env != "TEST" && dsn == "" {
		panic("database connection not set")
	}

	var gormDb = new(gorm.DB)
	var err = errors.New("")

	if env == "TEST" {
		gormDb, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
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

		gormDb, err = gorm.Open(mysql.Open(dsn), config)
	}

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db = &DB{Conn: gormDb}
	fmt.Println("Connected to Database")
}
