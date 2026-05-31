package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

func openTestDB(t *testing.T) *DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestProfileUpsertAndGet(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewProfileRepo(db)

	want := domain.Profile{UID: 42, FirstName: "Иван", Username: "ivan"}
	if err := repo.Upsert(ctx, want); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	got, err := repo.Get(ctx, 42)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.FirstName != "Иван" || got.Username != "ivan" {
		t.Fatalf("got %+v", got)
	}
	if got.Category != domain.CategoryUnknown {
		t.Fatalf("default category should be unknown, got %s", got.Category)
	}

	// Second upsert updates last_seen_at without erroring.
	if err := repo.Upsert(ctx, want); err != nil {
		t.Fatalf("re-upsert: %v", err)
	}
}

func TestProfileSetCategoryAndBlock(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewProfileRepo(db)
	_ = repo.Upsert(ctx, domain.Profile{UID: 1, FirstName: "X"})

	if err := repo.SetCategory(ctx, 1, domain.CategoryHR); err != nil {
		t.Fatalf("set cat: %v", err)
	}
	if err := repo.SetBlocked(ctx, 1, true); err != nil {
		t.Fatalf("set blocked: %v", err)
	}
	p, _ := repo.Get(ctx, 1)
	if p.Category != domain.CategoryHR || !p.IsBlocked {
		t.Fatalf("bad state: %+v", p)
	}
}

func TestProfileListFilter(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewProfileRepo(db)
	_ = repo.Upsert(ctx, domain.Profile{UID: 1, FirstName: "A"})
	_ = repo.Upsert(ctx, domain.Profile{UID: 2, FirstName: "B"})
	_ = repo.SetCategory(ctx, 1, domain.CategoryHR)
	_ = repo.SetCategory(ctx, 2, domain.CategoryCoop)

	hr, err := repo.List(ctx, ListFilter{Category: domain.CategoryHR})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(hr) != 1 || hr[0].UID != 1 {
		t.Fatalf("expected only uid=1 in HR, got %+v", hr)
	}
}

func TestMessageAppendAndList(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	prepo := NewProfileRepo(db)
	mrepo := NewMessageRepo(db)
	_ = prepo.Upsert(ctx, domain.Profile{UID: 7})

	if _, err := mrepo.Append(ctx, domain.Message{UID: 7, Direction: domain.DirectionIn, Text: "hi"}); err != nil {
		t.Fatalf("append in: %v", err)
	}
	if _, err := mrepo.Append(ctx, domain.Message{UID: 7, Direction: domain.DirectionOut, Text: "hey"}); err != nil {
		t.Fatalf("append out: %v", err)
	}
	msgs, err := mrepo.ListByUID(ctx, 7, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("want 2 msgs, got %d", len(msgs))
	}
	if msgs[0].Direction != domain.DirectionIn || msgs[1].Direction != domain.DirectionOut {
		t.Fatalf("wrong order: %+v", msgs)
	}
}

func TestStateSaveGetDelete(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	prepo := NewProfileRepo(db)
	srepo := NewStateRepo(db)
	_ = prepo.Upsert(ctx, domain.Profile{UID: 9})

	if _, err := srepo.Get(ctx, 9); err == nil {
		t.Fatal("expected ErrNotFound on first get")
	}
	s := domain.ConversationState{UID: 9, StateID: domain.StateAskCategory, Data: domain.StateData{"k": "v"}}
	if err := srepo.Save(ctx, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := srepo.Get(ctx, 9)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.StateID != domain.StateAskCategory || got.Data["k"] != "v" {
		t.Fatalf("bad state: %+v", got)
	}
	if err := srepo.Delete(ctx, 9); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := srepo.Get(ctx, 9); err == nil {
		t.Fatal("expected ErrNotFound after delete")
	}
}
