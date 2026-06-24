CREATE TABLE IF NOT EXISTS catalog_items (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(100) NOT NULL,
    slug         VARCHAR(50)  UNIQUE NOT NULL,
    base_price   NUMERIC(10,2) NOT NULL,
    is_active    BOOLEAN      NOT NULL DEFAULT true,
    sort_order   INTEGER      NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  DEFAULT NOW()
);

-- Seed catalog items matching WashPoint design
INSERT INTO catalog_items (name, slug, base_price, sort_order) VALUES
    ('Shirt',       'shirt',     1.50,  1),
    ('T-Shirt',     'tshirt',    1.20,  2),
    ('Trousers',    'trousers',  2.50,  3),
    ('Jeans',       'jeans',     2.50,  4),
    ('Dress',       'dress',     4.00,  5),
    ('Suit (2pc)',  'suit',      8.00,  6),
    ('Jacket',      'jacket',    5.00,  7),
    ('Bedsheet',    'bedsheet',  4.00,  8),
    ('Duvet',       'duvet',     9.00,  9),
    ('Towel',       'towel',     1.00, 10),
    ('Curtain',     'curtain',   6.00, 11),
    ('Saree',       'saree',     5.00, 12)
ON CONFLICT (slug) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_catalog_items_is_active ON catalog_items(is_active);
