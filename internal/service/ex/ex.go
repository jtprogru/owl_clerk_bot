package ex

import (
	"context"
)

type SM struct {
	store *StateStorage
}

func NewSM() *SM {
	return &SM{store: NewStateStorage()}
}

type answer struct {
	msg     string
	buttons []string
}

func (a answer) GetMessage() string {
	return a.msg
}
func (a answer) GetKeyboard() []string {
	return a.buttons
}

type Profile interface {
	GetUID() int64
	GetFirstName() string
	GetLastName() string
	GetUsername() string
}

type Message interface {
	GetMessage() string
}

type StateStorage struct {
	s map[int64]int
}

func (s *StateStorage) GetState(userId int64) State {
	stateId := s.s[userId]
	switch stateId {
	case 0:
		return StartState{}
	case 1:
		return WelcomeState{}
	case 2:
		return HRState{}
	case 3:
		return NonHRState{}
	default:
		return WelcomeState{}
	}
}

func (s *StateStorage) SetState(userId int64, stateId int) {
	s.s[userId] = stateId
}

func NewStateStorage() *StateStorage {
	return &StateStorage{s: map[int64]int{}}
}

func (s SM) SaveOrUpdateState(ctx context.Context, p Profile, m Message) (answer, error) {
	state := s.store.GetState(p.GetUID())
	state = state.SelectNextState(p, m)
	a := state.GenerateAnswer()
	s.store.SetState(p.GetUID(), state.GetID())
	return a, nil
}

type State interface {
	SelectNextState(p Profile, m Message) State
	GenerateAnswer() answer
	GetID() int
}

type StartState struct{}

func (s StartState) SelectNextState(p Profile, m Message) State {
	return WelcomeState{}
}

func (s StartState) GenerateAnswer() answer {
	return answer{}
}
func (s StartState) GetID() int {
	return 0
}

type WelcomeState struct{}

func (s WelcomeState) SelectNextState(p Profile, m Message) State {
	switch m.GetMessage() {
	case "HR":
		return HRState{}
	case "NonHR":
		return NonHRState{}
	default:
		return WelcomeState{}
	}
}

func (s WelcomeState) GenerateAnswer() answer {
	return answer{
		msg: "Who are you?",
		buttons: []string{
			"HR",
			"NonHR",
		},
	}
}

func (s WelcomeState) GetID() int {
	return 1
}

type HRState struct{}

func (s HRState) SelectNextState(p Profile, m Message) State {
	return HRState{}
}

func (s HRState) GenerateAnswer() answer {
	return answer{
		msg: "Ты HR",
		buttons: []string{
			"OK",
			"NOK",
		},
	}
}
func (s HRState) GetID() int {
	return 2
}

type NonHRState struct{}

func (s NonHRState) SelectNextState(p Profile, m Message) State {
	return NonHRState{}
}

func (s NonHRState) GenerateAnswer() answer {
	return answer{
		msg: "Ты NonHR",
		buttons: []string{
			"OK",
			"NOK",
		},
	}
}

func (s NonHRState) GetID() int {
	return 3
}
