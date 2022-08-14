package main

import (
    "context"
    "github.com/jtprogru/owl_clerk_bot/internal/service"
    "github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
    "github.com/sirupsen/logrus"
    "os"
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

    logger.Info("bot is started...")
    botToken := os.Getenv("OWL_BOT_TOKEN")
    cfg := &tg.Config{
        BotToken: botToken,
        IsDebug:  false,
    }
    storer := &MockStorer{}
    service := service.NewService(storer, storer)
    tg := tg.NewTG(service, logger, cfg)

    tg.Run()
}
