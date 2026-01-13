-- 037_group_billing_mode.sql
-- 新增分组计费模式字段，支持次卡计费

-- 新增分组计费模式字段
-- billing_mode: 'balance'(默认,按token扣余额), 'subscription'(订阅模式), 'card'(次卡固定价格)
ALTER TABLE groups ADD COLUMN IF NOT EXISTS billing_mode VARCHAR(20) DEFAULT 'balance';

-- 新增分组默认次卡价格（当 billing_mode = 'card' 时使用）
ALTER TABLE groups ADD COLUMN IF NOT EXISTS default_card_price DECIMAL(20,8) DEFAULT NULL;

-- 迁移现有数据：subscription_type = 'subscription' 的分组设置 billing_mode = 'subscription'
UPDATE groups SET billing_mode = 'subscription' WHERE subscription_type = 'subscription';

-- 添加注释
COMMENT ON COLUMN groups.billing_mode IS '计费模式: balance(余额), subscription(订阅), card(次卡)';
COMMENT ON COLUMN groups.default_card_price IS '次卡模式默认单次价格(USD)';
