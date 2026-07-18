-- +goose Up
CREATE TABLE houses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id    UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    listing_key     TEXT NOT NULL,
    label           TEXT NOT NULL,
    pos_x           DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_y           DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_z           DOUBLE PRECISION NOT NULL DEFAULT 0,
    purchase_price  BIGINT NOT NULL CHECK (purchase_price >= 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT houses_character_listing_unique UNIQUE (character_id, listing_key)
);

CREATE INDEX idx_houses_character_id ON houses(character_id);

-- MVP: one owned house property total per character (enforced in app; listing unique helps too)

-- +goose Down
DROP TABLE IF EXISTS houses;
