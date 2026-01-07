-- 030: Add group model rates table for per-model rate configuration
-- This allows setting different rate multipliers for specific models within a group

-- Create group_model_rates table
CREATE TABLE IF NOT EXISTS group_model_rates (
    id              BIGSERIAL PRIMARY KEY,
    group_id        BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    model           VARCHAR(100) NOT NULL,
    rate_multiplier DECIMAL(10, 4) NOT NULL DEFAULT 1.0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(group_id, model)
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_group_model_rates_group_id ON group_model_rates(group_id);
CREATE INDEX IF NOT EXISTS idx_group_model_rates_model ON group_model_rates(model);

-- Add rate_source column to usage_logs to track where the rate came from
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS rate_source VARCHAR(20) DEFAULT 'group';

-- Add comment for documentation
COMMENT ON TABLE group_model_rates IS 'Per-model rate multiplier configuration for groups';
COMMENT ON COLUMN group_model_rates.model IS 'Model name for exact matching';
COMMENT ON COLUMN group_model_rates.rate_multiplier IS 'Rate multiplier for this specific model';
COMMENT ON COLUMN usage_logs.rate_source IS 'Source of rate: group (default) or model (per-model config)';
