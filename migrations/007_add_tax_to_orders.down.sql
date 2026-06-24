ALTER TABLE orders
    DROP COLUMN IF EXISTS subtotal,
    DROP COLUMN IF EXISTS tax_rate,
    DROP COLUMN IF EXISTS tax_amount;
