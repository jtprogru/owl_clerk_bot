package sm

import "github.com/jtprogru/owl_clerk_bot/internal/domain"

var registry = map[domain.StateID]State{
	domain.StateGreeting:        greetingState{},
	domain.StateAskCategory:     askCategoryState{},
	domain.StateAskHRCompany:    askHRCompanyState{},
	domain.StateAskHRRole:       askHRRoleState{},
	domain.StateAskCoopProject:  askCoopProjectState{},
	domain.StateAskFriendOrigin: askFriendOriginState{},
	domain.StateAskMiscFree:     askMiscFreeState{},
	domain.StateAskContact:      askContactState{},
	domain.StateDone:            doneState{},
}

func resolve(id domain.StateID) (State, bool) {
	s, ok := registry[id]
	return s, ok
}
