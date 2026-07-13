BEGIN;

-- Using the postgis
CREATE EXTENSION IF NOT EXISTS postgis;


CREATE TABLE IF NOT EXISTS driver_locations (
    driver_id  UUID PRIMARY KEY REFERENCES drivers(id) ON DELETE CASCADE,
    location   GEOGRAPHY(Point, 4326) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );


CREATE INDEX IF NOT EXISTS idx_driver_locations_spatial ON driver_locations USING gist(location);
CREATE INDEX IF NOT EXISTS idx_driver_locations_updated_at ON driver_locations(updated_at);

COMMIT;