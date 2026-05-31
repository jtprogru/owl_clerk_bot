package sm

import (
	"testing"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

func TestMachine_FreshConversation(t *testing.T) {
	m := New()
	out, err := m.Step(42, domain.ConversationState{}, "anything")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.State.StateID != domain.StateGreeting {
		t.Fatalf("want Greeting, got %s", out.State.StateID)
	}
	if out.Reply.Text == "" {
		t.Fatal("expected greeting text")
	}
	if out.JustDone {
		t.Fatal("JustDone should be false on first step")
	}
}

func TestMachine_HRBranch(t *testing.T) {
	m := New()

	steps := []struct {
		input   string
		wantID  domain.StateID
		wantKey string
		wantVal string
	}{
		{"hi", domain.StateAskCategory, "", ""},
		{"HR", domain.StateAskHRCompany, KeyCategory, "HR"},
		{"Acme", domain.StateAskHRRole, KeyCompany, "Acme"},
		{"Senior Go", domain.StateAskContact, KeyRole, "Senior Go"},
		{"+7 999", domain.StateDone, KeyContact, "+7 999"},
	}

	state := domain.ConversationState{StateID: domain.StateGreeting, Data: domain.StateData{}}
	for i, st := range steps {
		out, err := m.Step(42, state, st.input)
		if err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
		if out.State.StateID != st.wantID {
			t.Fatalf("step %d: want state %s, got %s", i, st.wantID, out.State.StateID)
		}
		if st.wantKey != "" {
			if got := out.State.Data[st.wantKey]; got != st.wantVal {
				t.Fatalf("step %d: want %s=%q, got %q", i, st.wantKey, st.wantVal, got)
			}
		}
		state = out.State
	}

	if state.StateID != domain.StateDone {
		t.Fatalf("final state must be Done, got %s", state.StateID)
	}
}

func TestMachine_DoneStaysSilent(t *testing.T) {
	m := New()
	state := domain.ConversationState{
		UID:     1,
		StateID: domain.StateDone,
		Data:    domain.StateData{KeyCategory: "HR"},
	}
	out, err := m.Step(1, state, "follow-up")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.State.StateID != domain.StateDone {
		t.Fatalf("want Done, got %s", out.State.StateID)
	}
	if out.Reply.Text != "" {
		t.Fatalf("expected silent reply, got %q", out.Reply.Text)
	}
	if out.JustDone {
		t.Fatal("JustDone must be false when already Done")
	}
}

func TestMachine_JustDoneFlagOnFirstTransition(t *testing.T) {
	m := New()
	state := domain.ConversationState{
		UID:     1,
		StateID: domain.StateAskContact,
		Data:    domain.StateData{KeyCategory: "Разное", KeyTopic: "вопрос"},
	}
	out, err := m.Step(1, state, "mail@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.JustDone {
		t.Fatal("expected JustDone=true on first transition into Done")
	}
	if out.State.StateID != domain.StateDone {
		t.Fatalf("want Done, got %s", out.State.StateID)
	}
	if out.Reply.Text == "" {
		t.Fatal("expected non-empty Done prompt on first arrival")
	}
}

func TestMachine_AskCategoryUnknownInputStays(t *testing.T) {
	m := New()
	state := domain.ConversationState{
		UID:     1,
		StateID: domain.StateAskCategory,
		Data:    domain.StateData{},
	}
	out, err := m.Step(1, state, "что-то непонятное")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.State.StateID != domain.StateAskCategory {
		t.Fatalf("want AskCategory (stay), got %s", out.State.StateID)
	}
}
