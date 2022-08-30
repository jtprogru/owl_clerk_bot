package tg

import (
	"context"
	"errors"
	"time"

	"github.com/jtprogru/owl_clerk_bot/internal/conf"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type TG struct {
	logger *logrus.Logger
	b      *tele.Bot
	sm     StateMachine
}

type Profile struct {
	uid                    int64
	fName, lName, username string
}

func (p Profile) GetUID() int64 {
	return p.uid
}

func (p Profile) GetFirstName() string {
	return p.fName
}

func (p Profile) GetLastName() string {
	return p.lName
}

func (p Profile) GetUsername() string {
	return p.username
}

type Message struct {
	msg string
}

func (m Message) GetMessage() string {
	return m.msg
}

type StateMachine interface {
	SaveOrUpdateState(ctx context.Context, p Profile, m Message) (Answer, error)
}

type Answer interface {
	GetMessage() string
	GetKeyboard() []string
}

func (tg *TG) handleMessage(c tele.Context) error {
	ctx := context.WithValue(context.Background(), "xrayID", uuid.New())
	log := tg.logger.WithContext(ctx)
	log.Debug("message recived")

	p := Profile{}
	m := Message{}

	p.uid = c.Sender().ID
	p.fName = c.Sender().FirstName
	p.lName = c.Sender().LastName
	p.username = c.Sender().Username
	m.msg = c.Message().Text

	a, err := tg.sm.SaveOrUpdateState(ctx, p, m)
	if err != nil {
		log.WithError(err).Error("cant save or update answer")
		return err
	}

	if msg := a.GetMessage(); len(msg) > 0 {
		c.Bot().Reply(c.Message(), msg)
	}

	if kb := a.GetKeyboard(); len(kb) > 0 {
		var btns []tele.Row
		// Universal markup builders.
		menu := &tele.ReplyMarkup{ResizeKeyboard: true}

		for _, tkb := range kb {
			btns = append(btns, menu.Row(menu.Text(tkb)))
		}
		menu.Reply(btns...)
		err := c.Send("menu", menu)
		if err != nil {
			log.WithError(err).Error("cant send menu")
			return err
		}
	}

	return nil
}

// onPing simple pinger.
func (tg *TG) onPing(c tele.Context) error {
	ctx := context.WithValue(context.Background(), "xrayID", uuid.New())
	log := tg.logger.WithContext(ctx)
	log.Info("ping-pong")
	if c == nil {
		log.WithError(errors.New("ErrNilContext"))
		return errors.New("ErrNilContext")
	}

	reply, err := c.Bot().Reply(c.Message(), "PONG")
	if err != nil {
		log.WithError(err)
		return err
	}
	log.WithFields(logrus.Fields{"command": "/ping", "from ": c.Sender().Username, "msg": reply.Text}).Info("ping-pong")

	return nil
}

func (tg *TG) Run() {

	ctx := context.WithValue(context.Background(), "xrayID", uuid.New())
	log := tg.logger.WithContext(ctx)
	tg.b.Handle(tele.OnText, tg.handleMessage)
	tg.b.Handle("/ping", tg.onPing)

	log.Info("bot is starting")

	tg.b.Start()
}

func NewTG(sm StateMachine, conf *conf.Config) *TG {
	var err error
	logger := conf.GetLogger()
	tg := &TG{
		logger: logger,
		sm:     sm,
	}
	tg.b, err = tele.NewBot(tele.Settings{
		Token:   conf.GetBotToken(),
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: conf.GetIsDebug(),
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
