-- 036_add_card_price_to_model_rates.sql
-- 在 group_model_rates 表新增次卡价格字段，用于次卡计费模式

-- 新增次卡价格字段
ALTER TABLE group_model_rates
ADD COLUMN IF NOT EXISTS card_price DECIMAL(20,8) DEFAULT NULL;

COMMENT ON COLUMN group_model_rates.card_price IS '次卡模式单次请求价格(USD)，NULL表示不支持次卡';

-- usage_logs 表新增次卡计费标记
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS is_card_billing BOOLEAN DEFAULT FALSE;

-- 创建索引便于查询次卡使用记录
CREATE INDEX IF NOT EXISTS idx_usage_logs_is_card_billing ON usage_logs(is_card_billing) WHERE is_card_billing = TRUE;
