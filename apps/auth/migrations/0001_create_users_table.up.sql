CREATE TABLE IF NOT EXISTS users(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ---
    username varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    salt varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    ---
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

