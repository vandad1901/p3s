CREATE TABLE IF NOT EXISTS post(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    blog_id int NOT NULL REFERENCES blog(id),
    ---
    title varchar(256) NOT NULL,
    content jsonb NOT NULL,    
    ---
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

