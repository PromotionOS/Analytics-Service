-- Analytics schema migration
-- Goose runs this on startup

-- +goose Up
CREATE TABLE IF NOT EXISTS analytics.campaign_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    baseline_sales_per_day DECIMAL(12,2) DEFAULT 0,
    actual_sales_per_day DECIMAL(12,2) DEFAULT 0,
    lift DECIMAL(12,2),
    lift_percentage DECIMAL(8,4),
    total_funding_cost DECIMAL(12,2) DEFAULT 0,
    incremental_margin DECIMAL(12,2),
    roi DECIMAL(8,4),
    budget_burn_percent DECIMAL(5,2) DEFAULT 0,
    burned_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    redemption_count INTEGER DEFAULT 0,
    budget_exhausted_emitted BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_campaign_tenant UNIQUE (campaign_id, tenant_id)
);

CREATE INDEX IF NOT EXISTS idx_metrics_tenant ON analytics.campaign_metrics(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS analytics.campaign_metrics;
