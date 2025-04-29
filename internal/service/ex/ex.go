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

//func (s *StateStorage) GetState(userId int64) State {
//	stateId := s.s[userId]
//	switch stateId {
//	case 0:
//		return StartState{}
//	case 1:
//		return WelcomeState{}
//	case 2:
//		return HRState{}
//	case 3:
//		return NonHRState{}
//	default:
//		return WelcomeState{}
//	}
//}

func (s *StateStorage) SetState(userId int64, stateId int) {
	s.s[userId] = stateId
}

func NewStateStorage() *StateStorage {
	return &StateStorage{s: map[int64]int{}}
}

func (s SM) SaveOrUpdateState(ctx context.Context, p Profile, m Message) (answer, error) {
	//state := s.store.GetState(p.GetUID())
	//state = state.SelectNextState(p, m)
	//a := state.GenerateAnswer()
	//s.store.SetState(p.GetUID(), state.GetID())
	//return a, nil
	return answer{}, nil
}

type State interface {
	SelectNextState(p Profile, m Message) State
	GenerateAnswer() answer
	GetID() int
}

//type StartState struct{}
//
//func (s StartState) SelectNextState(p Profile, m Message) State {
//	return WelcomeState{}
//}
//
//func (s StartState) GenerateAnswer() answer {
//	return answer{}
//}
//func (s StartState) GetID() int {
//	return 0
//}
//
//type WelcomeState struct{}
//
//func (s WelcomeState) SelectNextState(p Profile, m Message) State {
//	switch m.GetMessage() {
//	case "HR":
//		return HRState{}
//	case "NonHR":
//		return NonHRState{}
//	default:
//		return WelcomeState{}
//	}
//}
//
//func (s WelcomeState) GenerateAnswer() answer {
//	return answer{
//		msg: "Who are you?",
//		buttons: []string{
//			"HR",
//			"NonHR",
//		},
//	}
//}
//
//func (s WelcomeState) GetID() int {
//	return 1
//}
//
//type HRState struct{}
//
//func (s HRState) SelectNextState(p Profile, m Message) State {
//	return HRState{}
//}
//
//func (s HRState) GenerateAnswer() answer {
//	return answer{
//		msg: "Ты HR",
//		buttons: []string{
//			"OK",
//			"NOK",
//		},
//	}
//}
//func (s HRState) GetID() int {
//	return 2
//}
//
//type NonHRState struct{}
//
//func (s NonHRState) SelectNextState(p Profile, m Message) State {
//	return NonHRState{}
//}
//
//func (s NonHRState) GenerateAnswer() answer {
//	return answer{
//		msg: "Ты NonHR",
//		buttons: []string{
//			"OK",
//			"NOK",
//		},
//	}
//}
//
//func (s NonHRState) GetID() int {
//	return 3
//}

type StateTable struct {
	Id       int
	NextIds  []int
	Answer   string
	Buttons  []string
	Heandler string
}

type GeneratedState struct {
	Id       int
	Answer   answer
	SelectNS func(Profile, Message) int
}

var HeandlerTable = map[string]HeandlerFunc{
	"ButtonSelect": ButtonSelect,
}

type HeandlerFunc func([]int, []string, int) func(p Profile, m Message) int

func (g GeneratedState) SelectNextState(p Profile, m Message) State {
	return GenerateState(GetStateTableByID(g.SelectNS(p, m)))
}

func (g GeneratedState) GenerateAnswer() answer {
	return g.Answer
}

func (g GeneratedState) GetID() int {
	return g.Id
}

func GenerateState(st StateTable) State {
	return GeneratedState{Id: st.Id, Answer: answer{st.Answer, st.Buttons}, SelectNS: HeandlerTable[st.Heandler](st.NextIds, st.Buttons, st.Id)}
}

func ButtonSelect(nx []int, bt []string, cs int) func(p Profile, m Message) int {
	return func(p Profile, m Message) int {
		if len(nx) == 0 {
			return cs
		}
		if len(nx) == 1 {
			return nx[0]
		}
		if len(bt) != len(nx) {
			return cs
		}
		for idx, btn := range bt {
			if btn == m.GetMessage() {
				return nx[idx]
			}
		}
		return cs
	}
}

func GetStateTableByID(int) StateTable {
	return StateTable{}
}
