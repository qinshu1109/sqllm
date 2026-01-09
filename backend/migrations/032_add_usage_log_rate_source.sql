-- Add rate_source field to usage_logs table
-- This field tracks whether the rate multiplier came from 'group' or 'model' level
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS rate_source VARCHAR(20) DEFAULT 'group';
