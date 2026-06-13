CREATE SCHEMA IF NOT EXISTS mkt_ecommerce;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_users(
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    user_type VARCHAR(20) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    is_agree BOOLEAN NOT NULL DEFAULT false, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    phone_number VARCHAR(50),
    user_locale VARCHAR(10),
    user_country VARCHAR(50),
    user_address TEXT,
    email_consent BOOLEAN NOT NULL DEFAULT false,
    sms_consent BOOLEAN NOT NULL DEFAULT false,
    consent_src VARCHAR(50),
    consent_updated_at TIMESTAMPTZ
);

CREATE INDEX idx_mkt_users_created_at ON mkt_ecommerce.mkt_users(created_at DESC);
CREATE UNIQUE INDEX mkt_users_email_key ON mkt_ecommerce.mkt_users(email);
CREATE UNIQUE INDEX mkt_users_username_key ON mkt_ecommerce.mkt_users(username);
CREATE INDEX idx_mkt_users_user_type ON mkt_ecommerce.mkt_users(user_type);
CREATE INDEX idx_mkt_users_email_consent ON mkt_ecommerce.mkt_users(email_consent);
CREATE INDEX idx_mkt_users_sms_consent ON mkt_ecommerce.mkt_users(sms_consent);
CREATE INDEX idx_mkt_users_unverified ON mkt_ecommerce.mkt_users(is_verified) WHERE is_verified = false;

CREATE TABLE IF NOT EXISTS mkt_ecommerce.mkt_refresh_tokens(
    token_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES mkt_ecommerce.mkt_users(user_id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_revoked BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_mkt_refresh_tokens_user_id ON mkt_ecommerce.mkt_refresh_tokens(user_id);
CREATE INDEX idx_mkt_refresh_tokens_token ON mkt_ecommerce.mkt_refresh_tokens(token);
CREATE INDEX idx_mkt_refresh_tokens_revoked ON mkt_ecommerce.mkt_refresh_tokens(is_revoked) WHERE is_revoked = false;