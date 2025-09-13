-- init.sql
CREATE TABLE IF NOT EXISTS todo (
    id TEXT PRIMARY KEY,
    todo TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deadline TIMESTAMP NOT NULL,
    priority TEXT,
    completed_at TIMESTAMP,
    complete BOOLEAN NOT NULL
);
