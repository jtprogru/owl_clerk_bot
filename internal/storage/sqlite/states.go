package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

type StateRepo struct {
	db *DB
}

func NewStateRepo(db *DB) *StateRepo { return &StateRepo{db: db} }

func (r *StateRepo) Get(ctx context.Context, uid int64) (domain.ConversationState, error) {
	const q = `SELECT uid, state_id, data, updated_at FROM conversation_states WHERE uid = ?`
	var (
		s       domain.ConversationState
		rawData string
	)
	err := r.db.QueryRowContext(ctx, q, uid).Scan(&s.UID, &s.StateID, &rawData, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ConversationState{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.ConversationState{}, fmt.Errorf("get state: %w", err)
	}
	if rawData != "" {
		if err := json.Unmarshal([]byte(rawData), &s.Data); err != nil {
			return domain.ConversationState{}, fmt.Errorf("decode state data: %w", err)
		}
	}
	if s.Data == nil {
		s.Data = domain.StateData{}
	}
	return s, nil
}

func (r *StateRepo) Save(ctx context.Context, s domain.ConversationState) error {
	if s.Data == nil {
		s.Data = domain.StateData{}
	}
	raw, err := json.Marshal(s.Data)
	if err != nil {
		return fmt.Errorf("encode state data: %w", err)
	}
	const q = `
INSERT INTO conversation_states(uid, state_id, data, updated_at)
VALUES(?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(uid) DO UPDATE SET
    state_id   = excluded.state_id,
    data       = excluded.data,
    updated_at = CURRENT_TIMESTAMP;`
	if _, err := r.db.ExecContext(ctx, q, s.UID, string(s.StateID), string(raw)); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	return nil
}

func (r *StateRepo) Delete(ctx context.Context, uid int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM conversation_states WHERE uid = ?`, uid)
	if err != nil {
		return fmt.Errorf("delete state: %w", err)
	}
	return nil
}
