package logger

import (
	"github.com/spf13/viper"
	"testing"
)

func TestLoggerInit(t *testing.T) {

	if ClientLog.Out != nil {
		t.Errorf("Init() = %v; want nil", ClientLog.Out)
	}

	if ClientLog.Out != nil {
		t.Errorf("Init() = %v; want nil", ClientLog.Out)
	}

	viper.Set("APP_ENV", "TEST")

	Init()

	if ClientLog.Out == nil {
		t.Errorf("Init() = %v; want value", ClientLog.Out)
	}

	if ClientLog.Level.String() != "info" {
		t.Errorf("Init() = %v; want info", ClientLog.Level)
	}

}
