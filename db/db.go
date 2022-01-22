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
	"os"
	"path/filepath"
)

type DB struct {
	Conn *gorm.DB
}

func NewConnection() *DB {
	dsn := viper.GetString("DB_DSN")
	env := viper.GetString("APP_ENV")

	var db = new(gorm.DB)
	var err = errors.New("")

	pwd, err := os.Getwd()

	if env == "TEST" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {

		file := util.GetFile(filepath.Join(pwd, "/logs/", "db.log"))
		if err != nil || file == nil {
			// Handle error
			panic("could not open file")
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		newLogger := logger.New(
			log.New(file, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				//SlowThreshold:              time.Second,   // Slow SQL threshold
				//LogLevel:                   logger.Silent, // Log level
				IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
				//Colorful:                  false,          // Disable color
			},
		)

		config := &gorm.Config{
			Logger: newLogger,
		}
		//
		//if env == "PRODUCTION" {
		//	// TODO: Add db logger configs to redirect sql logs to file
		//			config = &gorm.Config{}
		//}

		db, err = gorm.Open(mysql.Open(dsn), config)
	}

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	//defer cancel()

	return &DB{Conn: db}

}
