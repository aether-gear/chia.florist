-- Backfill SLA timestamps for legacy orders.
--
-- Sets `confirmed_at` from `updated_at` (or `created_at`) and
-- `handling_expires_at` to 3 days later.
--
-- Only updates `confirmed` and `processing` orders with missing values,
-- making this migration idempotent.
UPDATE orders
SET
    confirmed_at = COALESCE(confirmed_at, updated_at, created_at),
    handling_expires_at = COALESCE(handling_expires_at, COALESCE(updated_at, created_at) + INTERVAL '3 days')
WHERE status IN ('confirmed', 'processing')
  AND (confirmed_at IS NULL OR handling_expires_at IS NULL);
