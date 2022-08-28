package sm

import (
	"context"
	"github.com/google/uuid"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
	"github.com/sirupsen/logrus"
	"strconv"
)

// SM is a simple structure for State Machine
type SM struct {
	logger *logrus.Logger
}

// SaveOrUpdateState must save state for new user or update current state in SM as State Machine
func (sm SM) SaveOrUpdateState(ctx context.Context, p tg.Profile, m tg.Message) (tg.Answer, error) {
	log := sm.logger.WithContext(ctx).WithFields(logrus.Fields{"xrayID": uuid.New()})

	userID := p.GetUID()
	username := p.GetUsername()
	msg := m.GetMessage()

	log.WithFields(logrus.Fields{"user_id": userID, "username": username}).Debug("SaveOrUpdateState is called")
	return newAnswer(msg, []string{strconv.FormatInt(userID, 10), username}), nil
}

// NewSM create new instance of SM
func NewSM(logger *logrus.Logger) SM {
	return SM{
		logger: logger,
	}
}
