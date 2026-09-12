CREATE TABLE IF NOT EXISTS media(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    media_key varchar(255) NOT NULL UNIQUE,
    ---
    ingest_status int NOT NULL,
    mime_type varchar(255) NOT NULL,
    size bigint NOT NULL,
    ---
    leased_at timestamptz NOT NULL,
    try_count int NOT NULL
);

