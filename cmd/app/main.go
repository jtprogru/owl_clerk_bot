package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jtprogru/owl_clerk_bot/internal/service"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
)

type MockStorer struct{}

func (m MockStorer) SaveOrUpdate(ctx context.Context, uid int64, fName, lName, username string) error {
	logrus.Info("SaveOrUpdate is called")
	return nil
}
func (m MockStorer) Save(ctx context.Context, uid int64, msg string) error {
	logrus.Info("Save is called")
	return nil
}
func (m MockStorer) GetMessagesByUID(ctx context.Context, uid int64) ([]string, error) {
	logrus.Info("GetMessagesByUID is called")
	return nil, nil
}

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
	storer := &MockStorer{}
	s := service.NewService(storer, storer)
	client := tg.NewTG(s, logger, cfg)

	client.Run()
	logger.Info("bot stopped...")
}
