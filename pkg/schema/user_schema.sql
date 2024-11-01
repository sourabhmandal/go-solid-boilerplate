CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    email text UNIQUE NOT NULL,
    bio  text
);
