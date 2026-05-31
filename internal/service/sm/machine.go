package sm

import (
	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

// Outcome is the result of processing one user message through the FSM.
type Outcome struct {
	Reply    Reply                    // what to send back; empty Text = stay silent
	State    domain.ConversationState // updated state (including captured data)
	JustDone bool                     // true when this step transitioned INTO Done
}

type Machine struct{}

func New() *Machine { return &Machine{} }

// Step advances the conversation by one user message.
//
// prev is the current state; pass a zero ConversationState (StateID == "") for
// a brand-new conversation — in that case the FSM emits the greeting prompt
// and does not consume the input.
func (m *Machine) Step(uid int64, prev domain.ConversationState, input string) (Outcome, error) {
	if prev.StateID == "" {
		next, ok := resolve(domain.StateGreeting)
		if !ok {
			return Outcome{}, domain.ErrUnknownState
		}
		return Outcome{
			Reply: next.Prompt(),
			State: domain.ConversationState{
				UID:     uid,
				StateID: next.ID(),
				Data:    domain.StateData{},
			},
		}, nil
	}

	cur, ok := resolve(prev.StateID)
	if !ok {
		return Outcome{}, domain.ErrUnknownState
	}

	data := prev.Data
	if data == nil {
		data = domain.StateData{}
	}
	if key := cur.CaptureKey(); key != "" {
		data[key] = input
	}

	nextID := cur.Next(input)
	next, ok := resolve(nextID)
	if !ok {
		return Outcome{}, domain.ErrUnknownState
	}

	out := Outcome{
		State: domain.ConversationState{
			UID:     uid,
			StateID: nextID,
			Data:    data,
		},
	}

	switch {
	case prev.StateID != domain.StateDone && nextID == domain.StateDone:
		// First time hitting Done — send confirmation and signal owner notification.
		out.Reply = next.Prompt()
		out.JustDone = true
	case nextID == domain.StateDone:
		// Already past Done — stay silent. Intake will still notify the owner.
		out.Reply = Reply{}
	default:
		out.Reply = next.Prompt()
	}

	return out, nil
}
