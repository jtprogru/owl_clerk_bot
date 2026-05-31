package tg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jtprogru/owl_clerk_bot/internal/service/notify"
	"github.com/jtprogru/owl_clerk_bot/internal/service/outbound"
	tele "gopkg.in/telebot.v3"
)

var (
	_ notify.OwnerSender   = (*Client)(nil)
	_ outbound.UserSender = (*Client)(nil)
)

// Client wraps *tele.Bot with the outbound calls needed by the service layer:
// it implements notify.OwnerSender and outbound.UserSender.
type Client struct {
	b   *tele.Bot
	log *slog.Logger
}

func NewClient(b *tele.Bot, log *slog.Logger) *Client {
	return &Client{b: b, log: log}
}

func (c *Client) Bot() *tele.Bot { return c.b }

// SendOwnerNotification posts a survey summary to the owner with inline
// action buttons (reply hint / block / spam / change category).
func (c *Client) SendOwnerNotification(_ context.Context, ownerID, targetUID int64, summary string) error {
	chat := &tele.Chat{ID: ownerID}
	_, err := c.b.Send(chat, summary, ownerActionsKeyboard(targetUID))
	if err != nil {
		return fmt.Errorf("send owner notification: %w", err)
	}
	return nil
}

// SendToUser delivers an outbound text initiated by the owner (from /reply or
// the web UI) and returns the resulting Telegram message ID.
func (c *Client) SendToUser(_ context.Context, uid int64, text string) (int64, error) {
	chat := &tele.Chat{ID: uid}
	msg, err := c.b.Send(chat, text)
	if err != nil {
		return 0, fmt.Errorf("send to user %d: %w", uid, err)
	}
	return int64(msg.ID), nil
}

// NewBot builds a *tele.Bot with sane defaults (long polling, OnError logger).
func NewBot(token string, debug bool, log *slog.Logger) (*tele.Bot, error) {
	return tele.NewBot(tele.Settings{
		Token:   token,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: debug,
		OnError: func(err error, _ tele.Context) {
			log.Error("bot error", "err", err)
		},
	})
}
