CREATE EXTENSION IF NOT EXISTS "pqcrypto";

CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_products(
    product_id      UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id        UUID          NOT NULL REFERENCES mkt_ecommerce.mkt_stores(store_id) ON DELETE CASCADE,
    product_name    VARCHAR(100)  NOT NULL,
    product_desc    TEXT,
    product_url     TEXT,
    product_pic     TEXT,
    price           NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    category        VARCHAR(50)   NOT NULL,
    rating          NUMERIC(3,2)  NOT NULL DEFAULT 0.0 CHECK (rating >= 0 AND rating <= 5),
    availability    INTEGER       NOT NULL DEFAULT 0 CHECK (availability >= 0),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ
);

CREATE INDEX idx_mkt_products_store_id ON mkt_ecommerce.mkt_products(store_id);
CREATE INDEX idx_mkt_products_store_category ON mkt_ecommerce.mkt_products(category);


