BEGIN;

CREATE TABLE IF NOT EXISTS ride_offers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ride_id    UUID NOT NULL REFERENCES rides(id) ON DELETE CASCADE,
    driver_id  UUID NOT NULL REFERENCES drivers(id),
    status     VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (ride_id, driver_id)
    );

ALTER TABLE ride_offers
    ADD CONSTRAINT ride_offers_status_check CHECK (
        status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'EXPIRED')
        );

CREATE INDEX IF NOT EXISTS idx_ride_offers_ride_id ON ride_offers(ride_id);
CREATE INDEX IF NOT EXISTS idx_ride_offers_driver_id ON ride_offers(driver_id);

COMMIT;