CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    api_key VARCHAR(255),
    api_key_expiry TIMESTAMP
);

CREATE TABLE IF NOT EXISTS peers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    public_key VARCHAR(44) NOT NULL,
    assigned_ip INET NOT NULL UNIQUE,
    status BOOLEAN DEFAULT FALSE,
    is_gateway BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    created_on TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requesting_peer_id UUID REFERENCES peers(id),
    target_peer_id UUID REFERENCES peers(id),
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    is_valid BOOLEAN DEFAULT TRUE
);

INSERT INTO users (id, username, password_hash) 
VALUES (
    uuid_generate_v4(), 
    'hades', 
    crypt('12345678', gen_salt('bf'))
)
ON CONFLICT (username) DO NOTHING;
