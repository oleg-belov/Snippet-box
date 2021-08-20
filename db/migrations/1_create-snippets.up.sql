-- Create a `snippets` sequence.
CREATE SEQUENCE seq_snippets
    INCREMENT BY 10
    MINVALUE 10000
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