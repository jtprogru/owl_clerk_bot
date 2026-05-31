CREATE TABLE IF NOT EXISTS profiles (
    uid           INTEGER PRIMARY KEY,
    first_name    TEXT NOT NULL DEFAULT '',
    last_name     TEXT NOT NULL DEFAULT '',
    username      TEXT NOT NULL DEFAULT '',
    category      TEXT NOT NULL DEFAULT 'unknown',
    is_blocked    INTEGER NOT NULL DEFAULT 0,
    notes         TEXT NOT NULL DEFAULT '',
    contact       TEXT NOT NULL DEFAULT '',
    first_seen_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    uid       INTEGER NOT NULL REFERENCES profiles(uid) ON DELETE CASCADE,
    direction TEXT    NOT NULL CHECK(direction IN ('in','out')),
    text      TEXT    NOT NULL,
    tg_msg_id INTEGER NOT NULL DEFAULT 0,
    at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_messages_uid_at ON messages(uid, at);

CREATE TABLE IF NOT EXISTS conversation_states (
    uid        INTEGER PRIMARY KEY REFERENCES profiles(uid) ON DELETE CASCADE,
    state_id   TEXT NOT NULL,
    data       TEXT NOT NULL DEFAULT '{}',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
