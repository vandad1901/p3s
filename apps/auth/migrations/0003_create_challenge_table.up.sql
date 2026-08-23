CREATE TABLE IF NOT EXISTS challenge(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id bigint NOT NULL REFERENCES user_t(id) ON DELETE CASCADE,
    ---
    secret_hash bytea NOT NULL,
    challenge_type varchar(64) NOT NULL,
    metadata jsonb NOT NULL,
    ---
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz NOT NULL
);

