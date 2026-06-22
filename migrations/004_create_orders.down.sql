DROP INDEX IF EXISTS idx_orders_received_at;
DROP INDEX IF EXISTS idx_orders_created_by;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_customer_id;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
