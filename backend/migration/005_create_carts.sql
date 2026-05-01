CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE IF NOT EXISTS mkt_ecommerce.mkt_carts(
    cart_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES mkt_ecommerce.mkt_users(user_id) ON DELETE CASCADE
);



CREATE IF NOT EXISTS mkt_ecommerce.mkt_cart_items(
    cart_item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID REFERENCES mkt_ecommerce.mkt_carts(cart_id) ON DELETE CASCADE,
    product_id REFERENCES mkt_ecommerce.mkt_products(product_id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1
);