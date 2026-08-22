CREATE TABLE IF NOT EXISTS post(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ---
    title varchar(256) NOT NULL,
    slug varchar(256) NOT NULL UNIQUE,
    status int NOT NULL,
    ---
    created_at timestamptz NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamptz NOT NULL,
    updated_by bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS post_block(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    post_id bigint NOT NULL REFERENCES post ON DELETE CASCADE,
    position int NOT NULL,
    ---
    block_type int NOT NULL,
    text_content text,
    media_content varchar(64),
    metadata jsonb NOT NULL,
    ---
    CONSTRAINT post_post_id_position_unique UNIQUE (post_id, position) DEFERRABLE INITIALLY DEFERRED
);

