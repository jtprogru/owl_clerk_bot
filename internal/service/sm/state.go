package sm

import "github.com/jtprogru/owl_clerk_bot/internal/domain"

// Reply describes what the bot says to the user when entering a state.
// Empty Text means "stay silent" — used after the conversation has ended.
type Reply struct {
	Text     string
	Keyboard []string
}

// State is a node of the conversation FSM.
type State interface {
	ID() domain.StateID
	// CaptureKey is the key in StateData under which the user's input should be
	// stored when transitioning out of this state. Empty means do not capture.
	CaptureKey() string
	// Prompt is what the bot says when entering this state.
	Prompt() Reply
	// Next picks the next state based on user input.
	Next(input string) domain.StateID
}
