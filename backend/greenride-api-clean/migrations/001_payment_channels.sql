-- Migration: Add MTN MoMo and Stripe Payment Channels
-- Run this after the tables are created by GORM AutoMigrate

-- ============================================
-- PHASE 1: Payment Channels Configuration
-- ============================================

-- MTN MoMo Direct API - Sandbox (for testing)
INSERT INTO t_payment_channels (account_id, name, channel_code, payment_methods, status, config, remark, created_at, updated_at)
VALUES (
    'momo_sandbox_001',
    'MTN MoMo Sandbox',
    'momo',
    '["momo"]',
    'active',
    '{
        "environment": "sandbox",
        "subscription_key": "YOUR_SANDBOX_SUBSCRIPTION_KEY",
        "api_user_id": "YOUR_SANDBOX_API_USER_ID",
        "api_key": "YOUR_SANDBOX_API_KEY",
        "target_environment": "sandbox",
        "currency": "EUR",
        "timeout": 30
    }',
    'MTN MoMo Sandbox for testing - Replace credentials with real ones from MTN Developer Portal',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    config = VALUES(config),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- MTN MoMo Direct API - Production (inactive by default)
INSERT INTO t_payment_channels (account_id, name, channel_code, payment_methods, status, config, remark, created_at, updated_at)
VALUES (
    'momo_prod_001',
    'MTN MoMo Rwanda Production',
    'momo',
    '["momo"]',
    'inactive',
    '{
        "environment": "production",
        "subscription_key": "YOUR_PRODUCTION_SUBSCRIPTION_KEY",
        "api_user_id": "YOUR_PRODUCTION_API_USER_ID",
        "api_key": "YOUR_PRODUCTION_API_KEY",
        "target_environment": "rwandacollection",
        "currency": "RWF",
        "timeout": 30
    }',
    'MTN MoMo Production Rwanda - Activate after adding real credentials',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    config = VALUES(config),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- Stripe - Test Mode
INSERT INTO t_payment_channels (account_id, name, channel_code, payment_methods, status, config, remark, created_at, updated_at)
VALUES (
    'stripe_test_001',
    'Stripe Test',
    'stripe',
    '["card", "visa", "master", "amex"]',
    'active',
    '{
        "secret_key": "sk_test_YOUR_STRIPE_TEST_SECRET_KEY",
        "publishable_key": "pk_test_YOUR_STRIPE_TEST_PUBLISHABLE_KEY",
        "webhook_secret": "whsec_YOUR_STRIPE_TEST_WEBHOOK_SECRET",
        "currency": "eur",
        "statement_descriptor": "GREENRIDE",
        "timeout": 30
    }',
    'Stripe Test Mode for development - Replace with real test keys from Stripe Dashboard',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    config = VALUES(config),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- Stripe - Live Mode (inactive by default)
INSERT INTO t_payment_channels (account_id, name, channel_code, payment_methods, status, config, remark, created_at, updated_at)
VALUES (
    'stripe_live_001',
    'Stripe Live',
    'stripe',
    '["card", "visa", "master", "amex"]',
    'inactive',
    '{
        "secret_key": "sk_live_YOUR_STRIPE_LIVE_SECRET_KEY",
        "publishable_key": "pk_live_YOUR_STRIPE_LIVE_PUBLISHABLE_KEY",
        "webhook_secret": "whsec_YOUR_STRIPE_LIVE_WEBHOOK_SECRET",
        "currency": "eur",
        "statement_descriptor": "GREENRIDE",
        "timeout": 30
    }',
    'Stripe Live Mode for production - Activate after adding real credentials',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    config = VALUES(config),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- Verify insertions
SELECT account_id, name, channel_code, status FROM t_payment_channels WHERE channel_code IN ('momo', 'stripe');
