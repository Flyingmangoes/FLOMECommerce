CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_admin(
    admin_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firstname   VARCHAR(50)  NOT NULL,
    lastname    VARCHAR(50)  NOT NULL,
    username    VARCHAR(50)  UNIQUE NOT NULL,
    work_email  VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    phone_number VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
); 

CREATE INDEX IF NOT EXISTS idx_mkt_admin_created_at ON mkt_ecommerce.mkt_admin(created_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS mkt_admin_email_key ON mkt_ecommerce.mkt_admin(work_email);
CREATE UNIQUE INDEX IF NOT EXISTS mkt_admin_phone_number_key ON mkt_ecommerce.mkt_admin(phone_number);
CREATE UNIQUE INDEX IF NOT EXISTS mkt_admin_username_key ON mkt_ecommerce.mkt_admin(username);