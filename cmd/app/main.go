package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02T15:04:05",
		FullTimestamp:   true,
	}
	logger.Out = os.Stdout
	logger.SetReportCaller(true)
	logger.Info("bot is started...")
}
