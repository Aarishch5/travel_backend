BEGIN;

ALTER TABLE ride_offers
    ADD CONSTRAINT ride_offers_status_check CHECK (
        status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'EXPIRED')
        );

ALTER TABLE rides
    ADD CONSTRAINT rides_status_check CHECK
        (
        status IN ('REQUESTED', 'ACCEPTED', 'REJECTED',
                   'CANCELLED', 'COMPLETED', 'NO_DRIVERS_FOUND','REACHED_AT_DESTINATION')
        );

ALTER TABLE drivers
    ADD CONSTRAINT drivers_status_check CHECK (status IN ('ONLINE', 'OFFLINE', 'ON_TRIP'));

COMMIT;