-- Create a `snippets` sequence.
CREATE SEQUENCE seq_snippets
    INCREMENT BY 1
    MINVALUE 1
    CACHE 1
    NO CYCLE;

-- Create a `snippets` table.
CREATE TABLE snippets (
    id BIGINT       NOT NULL DEFAULT nextval('seq_snippets') PRIMARY KEY,
    title           VARCHAR(100) NOT NULL,
    content         TEXT NOT NULL,
    created         TIMESTAMP NOT NULL,
    expires         TIMESTAMP NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets(created);

-- Create a `users` sequence.
CREATE SEQUENCE seq_users
    INCREMENT BY 1
    MINVALUE 1
    CACHE 1
    NO CYCLE;

-- Create a `users` table.
CREATE TABLE users (
    id BIGINT       NOT NULL DEFAULT nextval('seq_users') PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2018-12-23 17:25:22 UTC'
);