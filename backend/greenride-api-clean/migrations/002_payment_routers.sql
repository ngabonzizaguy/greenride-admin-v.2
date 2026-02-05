-- Migration: Configure Payment Routers with Priority-based Routing
-- MoMo/Stripe get priority 200 (primary), KPay gets priority 100 (fallback)

-- ============================================
-- PHASE 2: Payment Routers Configuration
-- ============================================

-- Generate unique router IDs using UUID format
-- Note: Adjust these UUIDs or use your own ID generation

-- ---------------------------------------------
-- MTN MoMo Routes (Priority 200 - Primary)
-- ---------------------------------------------

-- MoMo Sandbox for testing (RWF)
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_momo_sandbox_rwf',
    'MoMo Sandbox - RWF',
    'momo',
    'momo_sandbox_001',
    'momo',
    'RWF',
    100.00,           -- Min amount (100 RWF)
    10000000.00,      -- Max amount (10M RWF)
    200,              -- High priority (primary)
    'active',
    'RW',             -- Rwanda region
    'MoMo sandbox route for testing',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- MoMo Production for Rwanda (inactive until credentials are added)
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_momo_prod_rwf',
    'MoMo Production - RWF',
    'momo',
    'momo_prod_001',
    'momo',
    'RWF',
    100.00,           -- Min amount (100 RWF)
    10000000.00,      -- Max amount (10M RWF)
    200,              -- High priority (primary)
    'inactive',       -- Activate after adding credentials
    'RW',             -- Rwanda region
    'MoMo production route for Rwanda',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- ---------------------------------------------
-- Stripe Routes (Priority 200 - Primary)
-- ---------------------------------------------

-- Stripe Test for EUR card payments
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_stripe_test_eur_card',
    'Stripe Test - EUR Cards',
    'stripe',
    'stripe_test_001',
    'card',
    'EUR',
    1.00,             -- Min amount (1 EUR)
    10000.00,         -- Max amount (10K EUR)
    200,              -- High priority (primary)
    'active',
    '*',              -- All regions
    'Stripe test route for EUR card payments',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- Stripe Test for USD card payments
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_stripe_test_usd_card',
    'Stripe Test - USD Cards',
    'stripe',
    'stripe_test_001',
    'card',
    'USD',
    1.00,             -- Min amount (1 USD)
    10000.00,         -- Max amount (10K USD)
    200,              -- High priority (primary)
    'active',
    '*',              -- All regions
    'Stripe test route for USD card payments',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- Stripe Live for EUR (inactive until credentials are added)
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_stripe_live_eur_card',
    'Stripe Live - EUR Cards',
    'stripe',
    'stripe_live_001',
    'card',
    'EUR',
    1.00,
    50000.00,
    200,
    'inactive',
    '*',
    'Stripe production route for EUR cards',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- ---------------------------------------------
-- KPay Fallback Routes (Priority 100 - Fallback)
-- ---------------------------------------------

-- Check if KPay routes exist, if so update their priority to 100
-- Otherwise these serve as examples of fallback configuration

-- KPay as fallback for MoMo payments
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_kpay_momo_fallback',
    'KPay MoMo Fallback',
    'kpay',
    'kpay_prod_001',  -- Adjust to your existing KPay account_id
    'momo',
    'RWF',
    100.00,
    10000000.00,
    100,              -- Lower priority (fallback)
    'active',
    'RW',
    'KPay fallback for MoMo when direct MoMo fails',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- KPay as fallback for card payments
INSERT INTO t_payment_routers (router_id, name, channel_code, channel_account_id, payment_method, currency, min_amount, max_amount, priority, status, region, remark, created_at, updated_at)
VALUES (
    'router_kpay_card_fallback',
    'KPay Card Fallback',
    'kpay',
    'kpay_prod_001',  -- Adjust to your existing KPay account_id
    'card',
    'RWF',
    100.00,
    10000000.00,
    100,              -- Lower priority (fallback)
    'active',
    'RW',
    'KPay fallback for card payments when Stripe fails',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    priority = VALUES(priority),
    updated_at = UNIX_TIMESTAMP() * 1000;

-- ---------------------------------------------
-- Update existing KPay routes to fallback priority
-- ---------------------------------------------

-- Lower priority of any existing KPay routes to make them fallbacks
UPDATE t_payment_routers
SET priority = 100,
    remark = CONCAT(COALESCE(remark, ''), ' [Updated to fallback priority]'),
    updated_at = UNIX_TIMESTAMP() * 1000
WHERE channel_code = 'kpay'
  AND priority > 100
  AND router_id NOT LIKE 'router_kpay%fallback';

-- Verify routing configuration
SELECT
    router_id,
    name,
    channel_code,
    payment_method,
    currency,
    priority,
    status,
    region
FROM t_payment_routers
WHERE channel_code IN ('momo', 'stripe', 'kpay')
ORDER BY payment_method, priority DESC;
