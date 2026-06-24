DROP INDEX IF EXISTS idx_orders_due_at;
DROP INDEX IF EXISTS idx_orders_payment_status;
DROP INDEX IF EXISTS idx_orders_order_number;
DROP SEQUENCE IF EXISTS order_number_seq;

ALTER TABLE orders
    DROP COLUMN IF EXISTS due_at,
    DROP COLUMN IF EXISTS payment_method,
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS service_type,
    DROP COLUMN IF EXISTS order_number;

ALTER TYPE order_status RENAME VALUE 'ready' TO 'done';
