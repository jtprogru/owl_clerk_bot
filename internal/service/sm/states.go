package sm

import (
	"strings"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

const (
	KeyCategory = "category"
	KeyCompany  = "company"
	KeyRole     = "role"
	KeyProject  = "project"
	KeyOrigin   = "origin"
	KeyTopic    = "topic"
	KeyContact  = "contact"
)

var categoryButtons = []string{"HR", "Сотрудничество", "Дружба", "Разное"}

type greetingState struct{}

func (greetingState) ID() domain.StateID    { return domain.StateGreeting }
func (greetingState) CaptureKey() string    { return "" }
func (greetingState) Next(string) domain.StateID {
	return domain.StateAskCategory
}
func (greetingState) Prompt() Reply {
	return Reply{
		Text: "Привет! Я бот-секретарь Михаила (@jtprogru). Помогаю ему не пропускать важное.\n\n" +
			"Сейчас задам пару коротких вопросов и передам тебя дальше.",
	}
}

type askCategoryState struct{}

func (askCategoryState) ID() domain.StateID { return domain.StateAskCategory }
func (askCategoryState) CaptureKey() string { return KeyCategory }
func (askCategoryState) Prompt() Reply {
	return Reply{
		Text:     "По какому ты вопросу?",
		Keyboard: categoryButtons,
	}
}
func (askCategoryState) Next(input string) domain.StateID {
	switch strings.TrimSpace(input) {
	case "HR":
		return domain.StateAskHRCompany
	case "Сотрудничество":
		return domain.StateAskCoopProject
	case "Дружба":
		return domain.StateAskFriendOrigin
	case "Разное":
		return domain.StateAskMiscFree
	}
	return domain.StateAskCategory
}

type askHRCompanyState struct{}

func (askHRCompanyState) ID() domain.StateID            { return domain.StateAskHRCompany }
func (askHRCompanyState) CaptureKey() string            { return KeyCompany }
func (askHRCompanyState) Prompt() Reply                 { return Reply{Text: "Из какой ты компании?"} }
func (askHRCompanyState) Next(string) domain.StateID    { return domain.StateAskHRRole }

type askHRRoleState struct{}

func (askHRRoleState) ID() domain.StateID         { return domain.StateAskHRRole }
func (askHRRoleState) CaptureKey() string         { return KeyRole }
func (askHRRoleState) Prompt() Reply              { return Reply{Text: "На какую роль/позицию?"} }
func (askHRRoleState) Next(string) domain.StateID { return domain.StateAskContact }

type askCoopProjectState struct{}

func (askCoopProjectState) ID() domain.StateID         { return domain.StateAskCoopProject }
func (askCoopProjectState) CaptureKey() string         { return KeyProject }
func (askCoopProjectState) Prompt() Reply              { return Reply{Text: "Расскажи в двух словах о проекте или идее."} }
func (askCoopProjectState) Next(string) domain.StateID { return domain.StateAskContact }

type askFriendOriginState struct{}

func (askFriendOriginState) ID() domain.StateID         { return domain.StateAskFriendOrigin }
func (askFriendOriginState) CaptureKey() string         { return KeyOrigin }
func (askFriendOriginState) Prompt() Reply              { return Reply{Text: "Откуда вы знакомы? Где познакомились?"} }
func (askFriendOriginState) Next(string) domain.StateID { return domain.StateAskContact }

type askMiscFreeState struct{}

func (askMiscFreeState) ID() domain.StateID         { return domain.StateAskMiscFree }
func (askMiscFreeState) CaptureKey() string         { return KeyTopic }
func (askMiscFreeState) Prompt() Reply              { return Reply{Text: "О чём хочешь поговорить?"} }
func (askMiscFreeState) Next(string) domain.StateID { return domain.StateAskContact }

type askContactState struct{}

func (askContactState) ID() domain.StateID         { return domain.StateAskContact }
func (askContactState) CaptureKey() string         { return KeyContact }
func (askContactState) Prompt() Reply              { return Reply{Text: "Оставь, пожалуйста, удобный способ связи (телефон, почта или @username)."} }
func (askContactState) Next(string) domain.StateID { return domain.StateDone }

type doneState struct{}

func (doneState) ID() domain.StateID         { return domain.StateDone }
func (doneState) CaptureKey() string         { return "" }
func (doneState) Prompt() Reply              { return Reply{Text: "Спасибо! Передал Михаилу, он напишет, как сможет."} }
func (doneState) Next(string) domain.StateID { return domain.StateDone }
