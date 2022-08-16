package tg

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type TG struct {
	storer Storer
	logger *logrus.Logger
	b      *tele.Bot
}

type Storer interface {
	StoreMessage(ctx context.Context, uid int64, msg string) error
	StoreOrUpdateProfile(ctx context.Context, uid int64, fName, lName, username string) error
}

func (tg *TG) handleMessage(c tele.Context) error {
	ctx := context.WithValue(context.Background(), "xrayID", uuid.New())
	tg.logger.WithContext(ctx).Debug("message recived")
	uid := c.Sender().ID
	fName := c.Sender().FirstName
	lName := c.Sender().LastName
	username := c.Sender().Username

	if err := tg.storer.StoreOrUpdateProfile(ctx, uid, fName, lName, username); err != nil {
		tg.logger.WithContext(ctx).WithError(err).Error("failed StoreOrUpdateProfile")
		return err
	}

	if err := tg.storer.StoreMessage(ctx, uid, c.Message().Text); err != nil {
		tg.logger.WithContext(ctx).WithError(err).Error("failed StoreMessage")
		return err
	}

	return nil
}

func (tg *TG) Run() {
	tg.b.Handle(tele.OnText, tg.handleMessage)

	tg.logger.Infof("bot is starting")

	tg.b.Start()
}

func NewTG(storer Storer, logger *logrus.Logger, cfg *Config) *TG {
	var err error

	tg := &TG{
		logger: logger,
		storer: storer,
	}
	tg.b, err = tele.NewBot(tele.Settings{
		Token:   cfg.BotToken,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: cfg.IsDebug,
		OnError: func(err error, c tele.Context) {
			logger.WithError(err).Error("bot on error")
		},
	})
	if err != nil {
		logger.WithError(err).Fatal("cant create bot")
		return nil
	}
	return tg

}
