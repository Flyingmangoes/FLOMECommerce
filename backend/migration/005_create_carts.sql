CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_carts(
    cart_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES mkt_ecommerce.mkt_users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_cart_items(
    cart_item_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id         UUID NOT NULL REFERENCES mkt_ecommerce.mkt_carts(cart_id) ON DELETE CASCADE,
    product_id      UUID NOT NULL REFERENCES mkt_ecommerce.mkt_products(product_id) ON DELETE CASCADE,
    store_id        UUID NOT NULL REFERENCES mkt_ecommerce.mkt_stores(store_id) ON DELETE CASCADE,
    quantity        INTEGER  NOT NULL DEFAULT 1 CHECK (quantity > 0),

    UNIQUE(cart_id, product_id)
);

CREATE INDEX idx_mkt_carts_user_id ON mkt_ecommerce.mkt_carts(user_id);
CREATE INDEX idx_mkt_cart_items_cart_id ON mkt_ecommerce.mkt_cart_items(cart_id);
CREATE INDEX idx_mkt_cart_items_quantity ON mkt_ecommerce.mkt_cart_items(quantity);
CREATE INDEX idx_mkt_cart_items_store_id ON mkt_ecommerce.mkt_cart_items(store_id);
CREATE INDEX idx_mkt_cart_items_product ON mkt_ecommerce.mkt_cart_items(product_id);