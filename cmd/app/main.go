package main

import (
	"context"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"

	// "github.com/jtprogru/owl_clerk_bot/internal/service"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
)

type MockSM struct{}

type MockAnswer struct {
	msg string
	kb  []string
}

func (ma MockAnswer) GetMessage() string {
	return ma.msg
}

func (ma MockAnswer) GetKeyboard() []string {
	return ma.kb
}

func (msm MockSM) SaveOrUpdateState(ctx context.Context, p tg.Profile, m tg.Message) (tg.Answer, error) {
	userID := p.GetUID()
	username := p.GetUsername()
	msg := m.GetMessage()

	return MockAnswer{
		msg: msg,
		kb:  []string{strconv.FormatInt(userID, 10), username},
	}, nil
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
	// s := service.NewService(storer, storer)
	mocksm := MockSM{}
	client := tg.NewTG(mocksm, logger, cfg)

	client.Run()
	logger.Info("bot stopped...")
}
