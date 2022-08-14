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
		ForceQuote:      true,
	}
	logger.Out = os.Stdout
	logger.SetReportCaller(true)
	logger.Info("bot is started...")
	botToken := os.Getenv("OWL_BOT_TOKEN")
	serveHost := os.Getenv("OWL_SERVE_HOST")
	servePort := os.Getenv("OWL_SERVE_PORT")
	logger.Infof("serve on host %v:%v with token %v", serveHost, servePort, botToken)
}
