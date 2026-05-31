package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
)

type ProfileRepo struct {
	db *DB
}

func NewProfileRepo(db *DB) *ProfileRepo { return &ProfileRepo{db: db} }

func (r *ProfileRepo) Upsert(ctx context.Context, p domain.Profile) error {
	const q = `
INSERT INTO profiles(uid, first_name, last_name, username, last_seen_at)
VALUES(?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(uid) DO UPDATE SET
    first_name  = excluded.first_name,
    last_name   = excluded.last_name,
    username    = excluded.username,
    last_seen_at = CURRENT_TIMESTAMP;`
	_, err := r.db.ExecContext(ctx, q, p.UID, p.FirstName, p.LastName, p.Username)
	if err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}
	return nil
}

func (r *ProfileRepo) Get(ctx context.Context, uid int64) (domain.Profile, error) {
	const q = `
SELECT uid, first_name, last_name, username, category, is_blocked, notes, contact,
       first_seen_at, last_seen_at
FROM profiles WHERE uid = ?;`
	var p domain.Profile
	var blocked int
	err := r.db.QueryRowContext(ctx, q, uid).Scan(
		&p.UID, &p.FirstName, &p.LastName, &p.Username,
		&p.Category, &blocked, &p.Notes, &p.Contact,
		&p.FirstSeenAt, &p.LastSeenAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Profile{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Profile{}, fmt.Errorf("get profile: %w", err)
	}
	p.IsBlocked = blocked != 0
	return p, nil
}

func (r *ProfileRepo) SetCategory(ctx context.Context, uid int64, c domain.Category) error {
	_, err := r.db.ExecContext(ctx, `UPDATE profiles SET category = ? WHERE uid = ?`, string(c), uid)
	if err != nil {
		return fmt.Errorf("set category: %w", err)
	}
	return nil
}

func (r *ProfileRepo) SetBlocked(ctx context.Context, uid int64, blocked bool) error {
	v := 0
	if blocked {
		v = 1
	}
	_, err := r.db.ExecContext(ctx, `UPDATE profiles SET is_blocked = ? WHERE uid = ?`, v, uid)
	if err != nil {
		return fmt.Errorf("set blocked: %w", err)
	}
	return nil
}

func (r *ProfileRepo) SetContact(ctx context.Context, uid int64, contact string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE profiles SET contact = ? WHERE uid = ?`, contact, uid)
	if err != nil {
		return fmt.Errorf("set contact: %w", err)
	}
	return nil
}

type ListFilter struct {
	Category domain.Category
	Blocked  *bool
	Limit    int
}

func (r *ProfileRepo) List(ctx context.Context, f ListFilter) ([]domain.Profile, error) {
	q := `
SELECT uid, first_name, last_name, username, category, is_blocked, notes, contact,
       first_seen_at, last_seen_at
FROM profiles WHERE 1=1`
	args := []any{}
	if f.Category != "" {
		q += ` AND category = ?`
		args = append(args, string(f.Category))
	}
	if f.Blocked != nil {
		q += ` AND is_blocked = ?`
		if *f.Blocked {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	q += ` ORDER BY last_seen_at DESC`
	if f.Limit > 0 {
		q += ` LIMIT ?`
		args = append(args, f.Limit)
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []domain.Profile
	for rows.Next() {
		var p domain.Profile
		var blocked int
		if err := rows.Scan(
			&p.UID, &p.FirstName, &p.LastName, &p.Username,
			&p.Category, &blocked, &p.Notes, &p.Contact,
			&p.FirstSeenAt, &p.LastSeenAt,
		); err != nil {
			return nil, fmt.Errorf("scan profile: %w", err)
		}
		p.IsBlocked = blocked != 0
		out = append(out, p)
	}
	return out, rows.Err()
}

type CategoryCount struct {
	Category domain.Category
	Count    int
}

func (r *ProfileRepo) CountByCategory(ctx context.Context) ([]CategoryCount, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT category, COUNT(*) FROM profiles GROUP BY category`)
	if err != nil {
		return nil, fmt.Errorf("count by category: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []CategoryCount
	for rows.Next() {
		var c CategoryCount
		if err := rows.Scan(&c.Category, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
