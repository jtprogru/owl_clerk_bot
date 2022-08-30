package service

import (
	"context"
	"strconv"

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
