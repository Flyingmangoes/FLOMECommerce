CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE IF NOT EXISTS mkt_ecommerce.mkt_carts(
    cart_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES mkt_ecommerce.mkt_users(user_id) ON DELETE CASCADE
);

CREATE IF NOT EXISTS mkt_ecommerce.mkt_cart_items(
    cart_item_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id         UUID NOT NULL REFERENCES mkt_ecommerce.mkt_carts(cart_id) ON DELETE CASCADE,
    product_id      UUID NOT NULL REFERENCES mkt_ecommerce.mkt_products(product_id) ON DELETE CASCADE,
    quantity        INT  NOT NULL DEFAULT 1 CHECK (quantity > 0),

    UNIQUE(cart_id, product_id)
);

CREATE INDEX idx_mkt_cart_user_id ON mkt_ecommerce.mkt_carts(user_id);
CREATE INDEX idx_mkt_cart_items_cart_id ON mkt_ecommerce.mkt_cart_items(cart_id);
CREATE INDEX idx_mkt_cart_items_product ON mkt_ecommerce.mkt_carts(product_id);