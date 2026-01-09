-- Create group_model_rates table for per-model rate configuration
CREATE TABLE IF NOT EXISTS group_model_rates (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    model VARCHAR(100) NOT NULL,
    rate_multiplier DECIMAL(10,4) NOT NULL DEFAULT 1.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(group_id, model)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_group_model_rates_group_id ON group_model_rates(group_id);
CREATE INDEX IF NOT EXISTS idx_group_model_rates_model ON group_model_rates(model);
