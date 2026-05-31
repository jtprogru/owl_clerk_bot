package tg

import (
	"strconv"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	tele "gopkg.in/telebot.v3"
)

// userReplyKeyboard builds a one-time reply keyboard from a list of button
// labels. Returns nil if labels is empty.
func userReplyKeyboard(labels []string) *tele.ReplyMarkup {
	if len(labels) == 0 {
		return nil
	}
	m := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	rows := make([]tele.Row, 0, len(labels))
	for _, l := range labels {
		rows = append(rows, m.Row(m.Text(l)))
	}
	m.Reply(rows...)
	return m
}

func clearReplyKeyboard() *tele.ReplyMarkup {
	return &tele.ReplyMarkup{RemoveKeyboard: true}
}

// Inline button uniques. telebot prefixes callback data with "\f<unique>|".
const (
	btnBlock  = "owl_block"
	btnSpam   = "owl_spam"
	btnCatPick = "owl_cat_pick"
	btnCatSet  = "owl_cat_set"
)

func ownerActionsKeyboard(uid int64) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	uidStr := strconv.FormatInt(uid, 10)
	block := m.Data("🚫 Заблок", btnBlock, uidStr)
	spam := m.Data("🗑 Спам", btnSpam, uidStr)
	catPick := m.Data("🔁 Сменить категорию", btnCatPick, uidStr)
	m.Inline(
		m.Row(block, spam),
		m.Row(catPick),
	)
	return m
}

// categoryPickerKeyboard returns an inline kbd to pick a new category for uid.
func categoryPickerKeyboard(uid int64) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	uidStr := strconv.FormatInt(uid, 10)
	mk := func(label string, cat domain.Category) tele.Btn {
		return m.Data(label, btnCatSet, uidStr+":"+string(cat))
	}
	m.Inline(
		m.Row(mk("HR", domain.CategoryHR), mk("Сотр.", domain.CategoryCoop)),
		m.Row(mk("Дружба", domain.CategoryFriend), mk("Разное", domain.CategoryMisc)),
		m.Row(mk("Спам", domain.CategorySpam)),
	)
	return m
}
