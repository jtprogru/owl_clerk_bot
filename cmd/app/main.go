package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/jtprogru/owl_clerk_bot/internal/ex"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
)

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

type csm struct {
	sm ex.SM
}

func (s csm) SaveOrUpdateState(ctx context.Context, p tg.IProfile, m tg.IMessage) (tg.Answer, error) {
	return s.sm.SaveOrUpdateState(ctx, p, m)
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
	mocksm := csm{
		sm: ex.SM{},
	}
	client := tg.NewTG(mocksm, logger, cfg)

	client.Run()
	logger.Info("bot stopped...")
}
