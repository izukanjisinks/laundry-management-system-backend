CREATE SEQUENCE IF NOT EXISTS order_number_seq START 1001;

DO $$ BEGIN
    CREATE TYPE order_status AS ENUM ('received', 'washing', 'ready', 'picked_up');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS orders (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number   INTEGER      NOT NULL DEFAULT nextval('order_number_seq'),
    customer_id    UUID         NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    created_by     UUID         NOT NULL REFERENCES users(id),
    status         order_status NOT NULL DEFAULT 'received',
    service_type   VARCHAR(20)  NOT NULL DEFAULT 'wash_fold',
    items          JSONB        NOT NULL,
    total_price    NUMERIC(10, 2),
    payment_status VARCHAR(20)  NOT NULL DEFAULT 'unpaid',
    payment_method VARCHAR(20),
    notes          TEXT,
    due_at         TIMESTAMPTZ,
    received_at    TIMESTAMPTZ  DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  DEFAULT NOW(),
    picked_up_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_orders_order_number   ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id    ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_status         ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_by     ON orders(created_by);
CREATE INDEX IF NOT EXISTS idx_orders_received_at    ON orders(received_at);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status);
CREATE INDEX IF NOT EXISTS idx_orders_due_at         ON orders(due_at);
