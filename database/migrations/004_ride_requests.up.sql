BEGIN;

CREATE TABLE IF NOT EXISTS rides (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id     UUID NOT NULL REFERENCES riders(id),
    driver_id    UUID REFERENCES drivers(id),
    status       VARCHAR(30) NOT NULL DEFAULT 'REQUESTED',
    pickup_lat   DOUBLE PRECISION NOT NULL,
    pickup_lng   DOUBLE PRECISION NOT NULL,
    drop_lat     DOUBLE PRECISION NOT NULL,
    drop_lng     DOUBLE PRECISION NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted_at  TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
    );

ALTER TABLE rides
    ADD CONSTRAINT rides_status_check CHECK (
        status IN ('REQUESTED', 'ACCEPTED', 'REJECTED', 'CANCELLED', 'COMPLETED', 'NO_DRIVERS_FOUND')
        );

CREATE INDEX IF NOT EXISTS idx_rides_driver_id ON rides(driver_id);
CREATE INDEX IF NOT EXISTS idx_rides_rider_id ON rides(rider_id);

COMMIT;

