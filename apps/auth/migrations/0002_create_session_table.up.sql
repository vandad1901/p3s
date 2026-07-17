CREATE TABLE IF NOT EXISTS sessions(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ---
    refresh_token_hash char(43) NOT NULL,
    ip_address inet,
    user_agent varchar(256) NOT NULL,
    status smallint NOT NULL,
    ---
    issued_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);

CREATE UNIQUE INDEX idx_sessions_refresh_token_hash ON sessions(refresh_token_hash);

CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

