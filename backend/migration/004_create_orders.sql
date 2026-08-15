CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_orders (
    order_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id    UUID NOT NULL REFERENCES mkt_ecommerce.mkt_users(user_id) ON DELETE CASCADE,
    buyer_email VARCHAR(100) NOT NULL,
    price_total NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    location    VARCHAR(150) NOT NULL,
    status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_mkt_orders_buyer_id ON mkt_ecommerce.mkt_orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_mkt_orders_buyer_email ON mkt_ecommerce.mkt_orders(buyer_email);

CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_order_items (
    order_item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID NOT NULL REFERENCES mkt_ecommerce.mkt_orders(order_id) ON DELETE CASCADE,
    product_id    UUID NOT NULL REFERENCES mkt_ecommerce.mkt_products(product_id) ON DELETE CASCADE,
    quantity      INTEGER NOT NULL CHECK (quantity > 0),
    price         NUMERIC(10, 2) NOT NULL CHECK (price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_mkt_order_items_order_id ON mkt_ecommerce.mkt_order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_mkt_order_items_product_id ON mkt_ecommerce.mkt_order_items(product_id);