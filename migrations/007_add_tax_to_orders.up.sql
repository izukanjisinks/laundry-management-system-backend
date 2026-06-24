ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS subtotal   NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tax_rate   NUMERIC(5,4)  NOT NULL DEFAULT 0.075,
    ADD COLUMN IF NOT EXISTS tax_amount NUMERIC(10,2) NOT NULL DEFAULT 0;

-- Backfill existing rows: treat current total_price as subtotal
UPDATE orders SET
    subtotal   = total_price,
    tax_amount = ROUND(total_price * 0.075, 2),
    total_price = ROUND(total_price + (total_price * 0.075), 2)
WHERE subtotal = 0;
