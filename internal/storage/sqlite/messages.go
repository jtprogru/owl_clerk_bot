package sqlite

import (
	"context"
	"fmt"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo { return &MessageRepo{db: db} }

func (r *MessageRepo) Append(ctx context.Context, m domain.Message) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO messages(uid, direction, text, tg_msg_id) VALUES(?, ?, ?, ?)`,
		m.UID, string(m.Direction), m.Text, m.TgMsgID,
	)
	if err != nil {
		return 0, fmt.Errorf("append message: %w", err)
	}
	return res.LastInsertId()
}

func (r *MessageRepo) ListByUID(ctx context.Context, uid int64, limit int) ([]domain.Message, error) {
	q := `SELECT id, uid, direction, text, tg_msg_id, at FROM messages WHERE uid = ? ORDER BY at ASC`
	args := []any{uid}
	if limit > 0 {
		q += ` LIMIT ?`
		args = append(args, limit)
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.UID, &m.Direction, &m.Text, &m.TgMsgID, &m.At); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *MessageRepo) CountSince(ctx context.Context, sinceDays int) (int, error) {
	q := fmt.Sprintf(`SELECT COUNT(*) FROM messages WHERE at >= datetime('now', '-%d days')`, sinceDays)
	var n int
	if err := r.db.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return 0, fmt.Errorf("count since: %w", err)
	}
	return n, nil
}
