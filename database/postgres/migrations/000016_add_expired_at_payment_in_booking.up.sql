ALTER TABLE bookings
    ADD COLUMN IF NOT EXISTS payment_expired_at timestamp;