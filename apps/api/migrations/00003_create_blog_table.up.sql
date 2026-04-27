CREATE TABLE IF NOT EXISTS blog(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(255) NOT NULL,
    tenant_id int NOT NULL REFERENCES tenant(id)
);

