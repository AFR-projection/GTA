-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    display_name    TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE characters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    gender          TEXT NOT NULL CHECK (gender IN ('male', 'female')),
    skin_tone       SMALLINT NOT NULL DEFAULT 0 CHECK (skin_tone BETWEEN 0 AND 10),
    hair_style      SMALLINT NOT NULL DEFAULT 0 CHECK (hair_style BETWEEN 0 AND 50),
    face_preset     SMALLINT NOT NULL DEFAULT 0 CHECK (face_preset BETWEEN 0 AND 50),
    outfit_id       TEXT NOT NULL DEFAULT 'starter_01',
    cash            BIGINT NOT NULL DEFAULT 500 CHECK (cash >= 0),
    bank            BIGINT NOT NULL DEFAULT 0 CHECK (bank >= 0),
    pos_x           DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_y           DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_z           DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT characters_account_name_unique UNIQUE (account_id, name)
);

-- MVP: one character per account (enforced in app + helpful partial uniqueness later if needed)
CREATE INDEX idx_characters_account_id ON characters(account_id);

CREATE TABLE inventory_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id    UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    item_key        TEXT NOT NULL,
    quantity        INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT inventory_character_item_unique UNIQUE (character_id, item_key)
);

CREATE INDEX idx_inventory_character_id ON inventory_items(character_id);

CREATE TABLE transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id    UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    kind            TEXT NOT NULL,
    amount          BIGINT NOT NULL,
    balance_cash    BIGINT,
    balance_bank    BIGINT,
    meta            JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_character_id ON transactions(character_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS accounts;
