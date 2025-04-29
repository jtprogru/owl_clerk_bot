package sm

import "context"

type StateStorage struct {
	s map[int64]int
}

type SM struct {
	store *StateStorage
}

func NewStateStorage() *StateStorage {
	return &StateStorage{s: map[int64]int{}}
}

func NewSM() *SM {
	return &SM{store: NewStateStorage()}
}

type Profile interface {
	GetUID() int64
	GetFirstName() string
	GetLastName() string
	GetUsername() string
}

type answer struct {
	msg     string
	buttons map[int]string
}

type Message interface {
	GetMessage() string
}

type State interface {
	SelectNextState(p Profile, m Message) State
	GenerateAnswer() answer
	GetID() int
}

func (s *StateStorage) GetState(userId int64) State {
	stateId := s.s[userId]
	switch stateId {
	case 0:
		return nil
		//StartState{}
	case 1:
	//	return WelcomeState{}
	//case 2:
	//	return HRState{}
	//case 3:
	//	return NonHRState{}
	default:
		return nil
		//WelcomeState{}
	}
}

func (a answer) GetMessage() string {
	return a.msg
}
func (a answer) GetKeyboard() map[int]string {
	return a.buttons
}

func (s SM) SaveOrUpdateState(ctx context.Context, p Profile, m Message) (answer, error) {
	state := s.store.GetState(p.GetUID())
	state = state.SelectNextState(p, m)
	a := state.GenerateAnswer()
	s.store.SetState(p.GetUID(), state.GetID())
	return a, nil
}

func (s *StateStorage) SetState(userId int64, stateId int) {
	s.s[userId] = stateId
}

type DefState struct{}

func (s DefState) SelectNextState(p Profile, m Message) State {
	switch m.GetMessage() {
	case "HR":
		return HRState{}
	case "NonHR":
		return NonHRState{}
	default:
		return WelcomeState{}
	}
}

func (s DefState) GenerateAnswer() answer {
	return answer{
		msg: "Who are you?",
		buttons: []string{
			"HR",
			"NonHR",
		},
	}
}

func (s DefState) GetID() int {
	return 1
}
