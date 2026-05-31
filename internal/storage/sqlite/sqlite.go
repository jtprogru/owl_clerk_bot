package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	*sql.DB
}

func Open(path string) (*DB, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)", path)
	raw, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	raw.SetMaxOpenConns(1) // SQLite serializes writes anyway; keep things simple.
	if err := raw.Ping(); err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	db := &DB{DB: raw}
	if err := db.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) migrate() error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, n := range names {
		raw, err := fs.ReadFile(migrationsFS, "migrations/"+n)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", n, err)
		}
		if _, err := db.Exec(string(raw)); err != nil {
			return fmt.Errorf("apply migration %s: %w", n, err)
		}
	}
	return nil
}
