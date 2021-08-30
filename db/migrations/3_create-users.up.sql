-- Create a `users` sequence.
CREATE SEQUENCE seq_users
    INCREMENT BY 10
    MINVALUE 10000
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
