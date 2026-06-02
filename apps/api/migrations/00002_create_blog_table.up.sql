CREATE TABLE IF NOT EXISTS blog(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    author_id int NOT NULL REFERENCES user(id),
    ---
    profile_picture varchar(256),
    description varchar(512),
    title varchar(256) NOT NULL,
);

