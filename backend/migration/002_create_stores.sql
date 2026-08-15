CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_stores(
    store_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES mkt_ecommerce.mkt_users(user_id) ON DELETE CASCADE,
    store_name VARCHAR(100) UNIQUE NOT NULL,
    store_desc TEXT,
    store_pic TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    store_locale VARCHAR(10),
    store_country VARCHAR(50),
    store_address TEXT, 
    store_phone_number VARCHAR(30),
    store_support_email VARCHAR(100),
    store_instagram VARCHAR(100),
    store_tiktok VARCHAR(100),
    store_website TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS mkt_stores_store_name ON mkt_ecommerce.mkt_stores(store_name);
CREATE INDEX IF NOT EXISTS idx_mkt_stores_owner_id ON mkt_ecommerce.mkt_stores(owner_id);
CREATE INDEX IF NOT EXISTS idx_mkt_stores_is_active ON mkt_ecommerce.mkt_stores(is_active);



