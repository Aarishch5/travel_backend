BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS drivers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(100) NOT NULL,
    email          VARCHAR(150) UNIQUE NOT NULL,
    phone          VARCHAR(20)  UNIQUE NOT NULL,
    license_number VARCHAR(50)  UNIQUE NOT NULL,
    vehicle_model  VARCHAR(50),
    plate_number   VARCHAR(20),
    created_at     TIMESTAMPTZ  DEFAULT now(),
    archived_at TIMESTAMPTZ,
    password_hash TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'OFFLINE'
    );

CREATE TABLE IF NOT EXISTS riders (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(150) UNIQUE NOT NULL,
    phone      VARCHAR(20)  UNIQUE NOT NULL,
    created_at TIMESTAMPTZ  DEFAULT now(),
    archived_at TIMESTAMPTZ,
    password_hash TEXT NOT NULL DEFAULT ''
    );

CREATE TABLE IF NOT EXISTS driver_locations (
    driver_id  UUID PRIMARY KEY REFERENCES drivers(id) ON DELETE CASCADE,
    latitude   DOUBLE PRECISION NOT NULL,
    longitude  DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS rides (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id     UUID NOT NULL REFERENCES riders(id),
    driver_id    UUID REFERENCES drivers(id),
    status       VARCHAR(30) NOT NULL DEFAULT 'REQUESTED',
    pickup_lat   DOUBLE PRECISION NOT NULL,
    pickup_lng   DOUBLE PRECISION NOT NULL,
    drop_lat     DOUBLE PRECISION NOT NULL,
    drop_lng     DOUBLE PRECISION NOT NULL,
    fare DECIMAL(10,2) DEFAULT 0.00,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted_at  TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
    );


CREATE TABLE IF NOT EXISTS ride_offers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ride_id    UUID NOT NULL REFERENCES rides(id) ON DELETE CASCADE,
    driver_id  UUID NOT NULL REFERENCES drivers(id),
    status     VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (ride_id, driver_id)
    );


COMMIT;