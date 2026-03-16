-- migrations/003_create_bookings.sql
CREATE TABLE IF NOT EXISTS bookings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ride_id     UUID NOT NULL REFERENCES rides(id) ON DELETE CASCADE,
    rider_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seats       INT NOT NULL DEFAULT 1 CHECK (seats > 0),
    status      TEXT NOT NULL DEFAULT 'pending',
    total_price NUMERIC(10,2) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT unique_rider_per_ride UNIQUE (ride_id, rider_id),
    CONSTRAINT valid_booking_status CHECK (
        status IN ('pending','confirmed','cancelled','completed')
    )
);

CREATE INDEX IF NOT EXISTS idx_bookings_ride   ON bookings(ride_id);
CREATE INDEX IF NOT EXISTS idx_bookings_rider  ON bookings(rider_id);