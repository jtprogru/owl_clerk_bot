package tg

import (
	"log/slog"
	"runtime/debug"

	tele "gopkg.in/telebot.v3"
)

func ownerOnly(ownerID int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender() == nil || c.Sender().ID != ownerID {
				return c.Reply("Доступ запрещён.")
			}
			return next(c)
		}
	}
}

func recoverMW(log *slog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("panic in handler",
						"recover", r,
						"stack", string(debug.Stack()),
					)
					err = nil
				}
			}()
			return next(c)
		}
	}
}
