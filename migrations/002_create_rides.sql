-- migrations/002_create_rides.sql
CREATE TABLE IF NOT EXISTS rides (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    origin_city         TEXT NOT NULL,
    destination_city    TEXT NOT NULL,
    origin_address      TEXT NOT NULL,
    destination_address TEXT NOT NULL,
    departure_at        TIMESTAMPTZ NOT NULL,
    total_seats         INT NOT NULL CHECK (total_seats > 0),
    available_seats     INT NOT NULL CHECK (available_seats >= 0),
    price_per_seat      NUMERIC(10,2) NOT NULL CHECK (price_per_seat >= 0),
    notes               TEXT,
    status              TEXT NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT available_lte_total CHECK (available_seats <= total_seats),
    CONSTRAINT valid_status CHECK (status IN ('active','full','cancelled','completed'))
);

CREATE INDEX IF NOT EXISTS idx_rides_search
    ON rides(origin_city, destination_city, departure_at, status);

CREATE INDEX IF NOT EXISTS idx_rides_driver
    ON rides(driver_id);