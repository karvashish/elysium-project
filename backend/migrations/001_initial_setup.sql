CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)), 2) || '-8' || substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6)))),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    api_key TEXT,
    api_key_expiry TEXT
);

CREATE TABLE IF NOT EXISTS peers (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)), 2) || '-8' || substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6)))),
    public_key TEXT NOT NULL,
    assigned_ip TEXT NOT NULL UNIQUE,
    status TEXT,
    is_gateway INTEGER DEFAULT 0,
    metadata TEXT,
    created_on TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)), 2) || '-8' || substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6)))),
    requesting_peer_id TEXT REFERENCES peers(id),
    target_peer_id TEXT REFERENCES peers(id),
    token TEXT NOT NULL,
    created_at TEXT NOT NULL,
    expired_at TEXT NOT NULL,
    is_valid INTEGER DEFAULT 1
);

INSERT INTO users (id, username, password_hash) 
VALUES (
    lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)), 2) || '-8' || substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6))), 
    'hades', 
    (SELECT hex(randomblob(16)))
)
ON CONFLICT (username) DO NOTHING;
