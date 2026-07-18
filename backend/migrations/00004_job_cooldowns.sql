-- +goose Up
CREATE TABLE job_cooldowns (
    character_id    UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    job_key         TEXT NOT NULL,
    last_completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (character_id, job_key)
);

-- +goose Down
DROP TABLE IF EXISTS job_cooldowns;
