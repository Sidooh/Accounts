package db

import (
	"github.com/spf13/viper"
	"testing"
)

func TestNewConnection(t *testing.T) {
	viper.Set("APP_ENV", "TEST")

	Init()
	conn := Connection()

	_, err := conn.DB()

	if err != nil {
		t.Errorf("NewConnection() failed: %s", conn.Error)
	}
}
