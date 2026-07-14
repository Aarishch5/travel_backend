BEGIN;

CREATE TABLE IF NOT EXISTS driver_locations (
    driver_id  UUID PRIMARY KEY REFERENCES drivers(id) ON DELETE CASCADE,
    latitude   DOUBLE PRECISION NOT NULL,
    longitude  DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_driver_locations_updated_at ON driver_locations(updated_at);

COMMIT;