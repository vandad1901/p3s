CREATE TABLE IF NOT EXISTS blog_user(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username varchar(255) NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

