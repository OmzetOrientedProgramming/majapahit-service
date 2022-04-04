ALTER TABLE bookings
    DROP COLUMN IF EXISTS xendit_id,
    DROP COLUMN IF EXISTS invoices_url