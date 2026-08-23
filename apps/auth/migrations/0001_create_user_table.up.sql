CREATE TABLE IF NOT EXISTS user_t(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ---
    username varchar(256) NOT NULL UNIQUE,
    password_hash varchar(256) NOT NULL,
    salt varchar(256) NOT NULL,
    email varchar(256) NOT NULL UNIQUE,
    ---
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

