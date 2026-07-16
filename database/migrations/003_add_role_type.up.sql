BEGIN;

CREATE TYPE role_type AS ENUM (
    'PASSENGER',
    'DRIVER',
    'ADMIN'
);

COMMIT;