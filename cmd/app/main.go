package main

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/jtprogru/owl_clerk_bot/internal/sm"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
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

	botToken := os.Getenv("OWL_BOT_TOKEN")
	cfg := &tg.Config{
		BotToken: botToken,
		IsDebug:  false,
	}

	s := sm.NewSM(logger)
	client := tg.NewTG(s, logger, cfg)

	client.Run()
}
