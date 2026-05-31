package intake

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	"github.com/jtprogru/owl_clerk_bot/internal/service/sm"
)

type fakeProfiles struct {
	store map[int64]domain.Profile
}

func newFakeProfiles() *fakeProfiles { return &fakeProfiles{store: map[int64]domain.Profile{}} }

func (f *fakeProfiles) Upsert(_ context.Context, p domain.Profile) error {
	existing, ok := f.store[p.UID]
	if !ok {
		p.Category = domain.CategoryUnknown
		f.store[p.UID] = p
		return nil
	}
	existing.FirstName = p.FirstName
	existing.LastName = p.LastName
	existing.Username = p.Username
	f.store[p.UID] = existing
	return nil
}
func (f *fakeProfiles) Get(_ context.Context, uid int64) (domain.Profile, error) {
	p, ok := f.store[uid]
	if !ok {
		return domain.Profile{}, domain.ErrNotFound
	}
	return p, nil
}
func (f *fakeProfiles) SetCategory(_ context.Context, uid int64, c domain.Category) error {
	p := f.store[uid]
	p.Category = c
	f.store[uid] = p
	return nil
}
func (f *fakeProfiles) SetContact(_ context.Context, uid int64, c string) error {
	p := f.store[uid]
	p.Contact = c
	f.store[uid] = p
	return nil
}

type fakeMessages struct {
	msgs []domain.Message
}

func (f *fakeMessages) Append(_ context.Context, m domain.Message) (int64, error) {
	f.msgs = append(f.msgs, m)
	return int64(len(f.msgs)), nil
}

type fakeStates struct {
	store map[int64]domain.ConversationState
}

func newFakeStates() *fakeStates { return &fakeStates{store: map[int64]domain.ConversationState{}} }

func (f *fakeStates) Get(_ context.Context, uid int64) (domain.ConversationState, error) {
	s, ok := f.store[uid]
	if !ok {
		return domain.ConversationState{}, domain.ErrNotFound
	}
	return s, nil
}
func (f *fakeStates) Save(_ context.Context, s domain.ConversationState) error {
	f.store[s.UID] = s
	return nil
}

type fakeNotifier struct {
	calls []notifyCall
}

type notifyCall struct {
	UID  int64
	Data domain.StateData
	Last string
}

func (n *fakeNotifier) Notify(_ context.Context, p domain.Profile, d domain.StateData, last string) error {
	n.calls = append(n.calls, notifyCall{UID: p.UID, Data: d, Last: last})
	return nil
}

func newSvc() (*Service, *fakeProfiles, *fakeMessages, *fakeStates, *fakeNotifier) {
	p := newFakeProfiles()
	m := &fakeMessages{}
	s := newFakeStates()
	n := &fakeNotifier{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := New(p, m, s, sm.New(), n, log)
	return svc, p, m, s, n
}

func TestIntake_FreshUserGetsGreeting(t *testing.T) {
	svc, _, _, states, notifier := newSvc()
	resp, err := svc.Handle(context.Background(), Incoming{
		Profile: domain.Profile{UID: 1, FirstName: "Иван"},
		Text:    "hello",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Reply.Text == "" {
		t.Fatal("expected greeting reply")
	}
	if got := states.store[1].StateID; got != domain.StateGreeting {
		t.Fatalf("want Greeting, got %s", got)
	}
	if len(notifier.calls) != 0 {
		t.Fatal("notifier should not fire on greeting")
	}
}

func TestIntake_FullHRFlowNotifiesOwner(t *testing.T) {
	svc, profiles, _, _, notifier := newSvc()
	uid := int64(7)
	p := domain.Profile{UID: uid, FirstName: "X"}

	// Greeting consumes the first message without capture, AskCategory captures
	// the second, and so on — full HR flow takes 6 messages.
	steps := []string{"start", "anything", "HR", "Acme", "Senior Go", "+7 999 123-45-67"}
	for _, txt := range steps {
		if _, err := svc.Handle(context.Background(), Incoming{Profile: p, Text: txt}); err != nil {
			t.Fatalf("step %q: %v", txt, err)
		}
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 owner notification, got %d", len(notifier.calls))
	}
	if got := profiles.store[uid].Category; got != domain.CategoryHR {
		t.Fatalf("want HR, got %s", got)
	}
	if got := profiles.store[uid].Contact; got != "+7 999 123-45-67" {
		t.Fatalf("want contact set, got %q", got)
	}
}

func TestIntake_BlockedUserSilentlyDropped(t *testing.T) {
	svc, profiles, messages, _, notifier := newSvc()
	uid := int64(13)
	_ = profiles.Upsert(context.Background(), domain.Profile{UID: uid})
	blocked := profiles.store[uid]
	blocked.IsBlocked = true
	profiles.store[uid] = blocked

	resp, err := svc.Handle(context.Background(), Incoming{
		Profile: domain.Profile{UID: uid, FirstName: "Spam"},
		Text:    "buy now",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !resp.Blocked {
		t.Fatal("expected Blocked=true")
	}
	if resp.Reply.Text != "" {
		t.Fatalf("expected silent reply, got %q", resp.Reply.Text)
	}
	if len(notifier.calls) != 0 {
		t.Fatal("notifier must not fire for blocked")
	}
	// Inbound message is still logged for the audit trail.
	if len(messages.msgs) != 1 {
		t.Fatalf("inbound msg must be saved, got %d", len(messages.msgs))
	}
}

func TestIntake_PostDoneMessagesAlsoNotify(t *testing.T) {
	svc, _, _, states, notifier := newSvc()
	uid := int64(21)
	states.store[uid] = domain.ConversationState{
		UID:     uid,
		StateID: domain.StateDone,
		Data:    domain.StateData{sm.KeyCategory: "HR"},
	}

	if _, err := svc.Handle(context.Background(), Incoming{
		Profile: domain.Profile{UID: uid},
		Text:    "follow-up question",
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(notifier.calls) != 1 {
		t.Fatalf("expected notification for follow-up msg, got %d", len(notifier.calls))
	}
	if notifier.calls[0].Last != "follow-up question" {
		t.Fatalf("want last msg passed to notifier, got %q", notifier.calls[0].Last)
	}
}
