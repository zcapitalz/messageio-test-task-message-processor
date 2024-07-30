DO $$ BEGIN
    CREATE TYPE message_status AS ENUM('not-processed', 'processed');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE DOMAIN ksuid AS CHAR(27);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS messages (
    id ksuid PRIMARY KEY,
    text TEXT NOT NULL,
    status message_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    processed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS messages_processed_at ON messages(processed_at);