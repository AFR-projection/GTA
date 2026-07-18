-- +goose Up
CREATE TABLE vehicles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id    UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    listing_key     TEXT NOT NULL,
    label           TEXT NOT NULL,
    vehicle_type    TEXT NOT NULL CHECK (vehicle_type IN ('motorcycle', 'car')),
    fuel            DOUBLE PRECISION NOT NULL DEFAULT 100 CHECK (fuel >= 0),
    fuel_max        DOUBLE PRECISION NOT NULL DEFAULT 100 CHECK (fuel_max > 0),
    pos_x           DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_y           DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_z           DOUBLE PRECISION NOT NULL DEFAULT 0,
    purchase_price  BIGINT NOT NULL CHECK (purchase_price >= 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vehicles_character_id ON vehicles(character_id);

-- +goose Down
DROP TABLE IF EXISTS vehicles;
