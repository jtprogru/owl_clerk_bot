package tg

import (
	"regexp"
	"strconv"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	tele "gopkg.in/telebot.v3"
)

func profileFromCtx(c tele.Context) domain.Profile {
	sender := c.Sender()
	if sender == nil {
		return domain.Profile{}
	}
	return domain.Profile{
		UID:       sender.ID,
		FirstName: sender.FirstName,
		LastName:  sender.LastName,
		Username:  sender.Username,
	}
}

// uidTagRe extracts "#uid_NNN" hash tag we embed in owner notifications, so
// the owner can simply reply-to a notification message and have the bot route
// the answer to the correct user.
var uidTagRe = regexp.MustCompile(`#uid_(\d+)`)

func extractTargetUID(text string) (int64, bool) {
	m := uidTagRe.FindStringSubmatch(text)
	if len(m) < 2 {
		return 0, false
	}
	uid, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, false
	}
	return uid, true
}
