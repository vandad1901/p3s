CREATE TABLE IF NOT EXISTS media_upload(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    key varchar(64) NOT NULL UNIQUE,
    ---
    status int NOT NULL,
    ---
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
);

