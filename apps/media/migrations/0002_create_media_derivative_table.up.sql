CREATE TABLE IF NOT EXISTS media_derivative(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    media_id bigint NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    media_key varchar(255) NOT NULL,
    ---
    ext varchar(5) NOT NULL,
    width int NOT NULL
);

