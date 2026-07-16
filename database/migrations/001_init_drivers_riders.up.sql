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
    avg_rating     NUMERIC(3,2) DEFAULT 5.0,
    created_at     TIMESTAMPTZ  DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS riders (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(150) UNIQUE NOT NULL,
    phone      VARCHAR(20)  UNIQUE NOT NULL,
    created_at TIMESTAMPTZ  DEFAULT now()
    );


COMMIT;