package global

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"main/internal/settings"
)

var (
	Logger *logrus.Entry
	Config *settings.Config
)

func LogDbMessage(table string, action string, message string) {
	msg := fmt.Sprintf("table: %v, action: %v, error: %v", table, action, message)
	Logger.Error(msg)
}
