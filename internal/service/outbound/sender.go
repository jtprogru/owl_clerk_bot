package outbound

import (
	"context"
	"fmt"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

// UserSender sends an arbitrary text to a Telegram user UID. Implemented by
// the Telegram transport.
type UserSender interface {
	SendToUser(ctx context.Context, uid int64, text string) (tgMsgID int64, err error)
}

// MessageStore is the slice of the storage layer we need to log outbound
// messages.
type MessageStore interface {
	Append(ctx context.Context, m domain.Message) (int64, error)
}

// Sender orchestrates an "out-of-band" reply: a message initiated by the
// owner (from a /reply command or the web UI) rather than by the FSM.
type Sender struct {
	sender   UserSender
	messages MessageStore
}

func New(s UserSender, m MessageStore) *Sender {
	return &Sender{sender: s, messages: m}
}

func (s *Sender) Reply(ctx context.Context, uid int64, text string) error {
	if text == "" {
		return fmt.Errorf("outbound: %w: empty text", domain.ErrInvalidInput)
	}
	tgID, err := s.sender.SendToUser(ctx, uid, text)
	if err != nil {
		return fmt.Errorf("send to user %d: %w", uid, err)
	}
	if _, err := s.messages.Append(ctx, domain.Message{
		UID:       uid,
		Direction: domain.DirectionOut,
		Text:      text,
		TgMsgID:   tgID,
	}); err != nil {
		return fmt.Errorf("log outbound: %w", err)
	}
	return nil
}
