package domain

import "time"

type StateID string

const (
	StateGreeting        StateID = "greeting"
	StateAskCategory     StateID = "ask_category"
	StateAskHRCompany    StateID = "ask_hr_company"
	StateAskHRRole       StateID = "ask_hr_role"
	StateAskCoopProject  StateID = "ask_coop_project"
	StateAskFriendOrigin StateID = "ask_friend_origin"
	StateAskMiscFree     StateID = "ask_misc_free"
	StateAskContact      StateID = "ask_contact"
	StateDone            StateID = "done"
)

type StateData map[string]string

type ConversationState struct {
	UID       int64
	StateID   StateID
	Data      StateData
	UpdatedAt time.Time
}

func (s *ConversationState) Set(key, value string) {
	if s.Data == nil {
		s.Data = StateData{}
	}
	s.Data[key] = value
}
