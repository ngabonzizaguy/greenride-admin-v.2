-- Migration: Activate MTN MoMo Production for Rwanda
-- Credentials from MTN Partner Portal (partner.mtn.co.rw) for GREENRIDE Ltd

-- ============================================
-- STEP 1: Update MoMo Production Channel with real credentials
-- ============================================
UPDATE t_payment_channels
SET config = '{
    "environment": "production",
    "subscription_key": "68866ad7ed7d4015883e0dd5f594b5eb",
    "api_user_id": "bfd36e1f-2587-46db-9700-062e49684c0f",
    "api_key": "40c71cbc476b46ec8ba612f239967046",
    "target_environment": "mtnrwanda",
    "currency": "RWF",
    "timeout": 30
}',
status = 'active',
remark = 'MTN MoMo Production Rwanda - Activated with live credentials Feb 2026',
updated_at = UNIX_TIMESTAMP() * 1000
WHERE account_id = 'momo_prod_001';

-- ============================================
-- STEP 2: Activate production payment router
-- ============================================
UPDATE t_payment_routers
SET status = 'active',
    updated_at = UNIX_TIMESTAMP() * 1000
WHERE router_id = 'router_momo_prod_rwf';

-- ============================================
-- STEP 3: Deactivate sandbox route (no longer needed)
-- ============================================
UPDATE t_payment_routers
SET status = 'inactive',
    updated_at = UNIX_TIMESTAMP() * 1000
WHERE router_id = 'router_momo_sandbox_rwf';

-- ============================================
-- VERIFY: Check the activated configuration
-- ============================================
SELECT account_id, name, channel_code, status, config
FROM t_payment_channels
WHERE account_id = 'momo_prod_001';

SELECT router_id, name, channel_code, payment_method, currency, priority, status
FROM t_payment_routers
WHERE channel_code = 'momo'
ORDER BY priority DESC;
