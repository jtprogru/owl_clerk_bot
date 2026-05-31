package intake

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	"github.com/jtprogru/owl_clerk_bot/internal/service/sm"
)

type ProfileStore interface {
	Upsert(ctx context.Context, p domain.Profile) error
	Get(ctx context.Context, uid int64) (domain.Profile, error)
	SetCategory(ctx context.Context, uid int64, c domain.Category) error
	SetContact(ctx context.Context, uid int64, contact string) error
}

type MessageStore interface {
	Append(ctx context.Context, m domain.Message) (int64, error)
}

type StateStore interface {
	Get(ctx context.Context, uid int64) (domain.ConversationState, error)
	Save(ctx context.Context, s domain.ConversationState) error
}

type Machine interface {
	Step(uid int64, prev domain.ConversationState, input string) (sm.Outcome, error)
}

type Notifier interface {
	Notify(ctx context.Context, p domain.Profile, data domain.StateData, lastMsg string) error
}

type Incoming struct {
	Profile domain.Profile
	Text    string
	TgMsgID int64
}

type Response struct {
	Reply   sm.Reply
	Blocked bool
}

type Service struct {
	profiles ProfileStore
	messages MessageStore
	states   StateStore
	machine  Machine
	notifier Notifier
	log      *slog.Logger
}

func New(p ProfileStore, m MessageStore, st StateStore, machine Machine, n Notifier, log *slog.Logger) *Service {
	return &Service{
		profiles: p, messages: m, states: st,
		machine: machine, notifier: n, log: log,
	}
}

func (s *Service) Handle(ctx context.Context, in Incoming) (Response, error) {
	if err := s.profiles.Upsert(ctx, in.Profile); err != nil {
		return Response{}, err
	}
	p, err := s.profiles.Get(ctx, in.Profile.UID)
	if err != nil {
		return Response{}, err
	}

	if _, err := s.messages.Append(ctx, domain.Message{
		UID:       p.UID,
		Direction: domain.DirectionIn,
		Text:      in.Text,
		TgMsgID:   in.TgMsgID,
	}); err != nil {
		return Response{}, err
	}

	if p.IsBlocked {
		s.log.Info("blocked user, ignoring", "uid", p.UID)
		return Response{Blocked: true}, nil
	}

	prev, err := s.states.Get(ctx, p.UID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return Response{}, err
	}

	out, err := s.machine.Step(p.UID, prev, in.Text)
	if err != nil {
		return Response{}, err
	}
	if err := s.states.Save(ctx, out.State); err != nil {
		return Response{}, err
	}

	if out.Reply.Text != "" {
		if _, err := s.messages.Append(ctx, domain.Message{
			UID:       p.UID,
			Direction: domain.DirectionOut,
			Text:      out.Reply.Text,
		}); err != nil {
			return Response{}, err
		}
	}

	if out.JustDone {
		if cat, ok := domain.ParseCategoryFromTitle(out.State.Data[sm.KeyCategory]); ok {
			if err := s.profiles.SetCategory(ctx, p.UID, cat); err != nil {
				s.log.Error("set category", "uid", p.UID, "err", err)
			} else {
				p.Category = cat
			}
		}
		if contact := out.State.Data[sm.KeyContact]; contact != "" {
			if err := s.profiles.SetContact(ctx, p.UID, contact); err != nil {
				s.log.Error("set contact", "uid", p.UID, "err", err)
			} else {
				p.Contact = contact
			}
		}
	}

	// Notify owner on:
	//  - just completed the survey
	//  - any subsequent message after Done
	shouldNotify := out.JustDone ||
		(prev.StateID == domain.StateDone && out.State.StateID == domain.StateDone)
	if shouldNotify {
		if err := s.notifier.Notify(ctx, p, out.State.Data, in.Text); err != nil {
			// Notification failures are non-fatal — we don't want to break the
			// user-facing reply because of a transient owner-side error.
			s.log.Error("notify owner", "uid", p.UID, "err", err)
		}
	}

	return Response{Reply: out.Reply}, nil
}
