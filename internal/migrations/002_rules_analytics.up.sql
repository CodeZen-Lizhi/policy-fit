CREATE TABLE IF NOT EXISTS rule_config_version (
    id BIGSERIAL PRIMARY KEY,
    version VARCHAR(64) NOT NULL UNIQUE,
    changelog TEXT NOT NULL,
    content JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_gray BOOLEAN NOT NULL DEFAULT FALSE,
    created_by BIGINT REFERENCES user_account(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rule_config_active ON rule_config_version(is_active);
CREATE INDEX IF NOT EXISTS idx_rule_config_created_at ON rule_config_version(created_at);

CREATE TABLE IF NOT EXISTS analytics_event (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES user_account(id) ON DELETE SET NULL,
    task_id BIGINT REFERENCES analysis_task(id) ON DELETE SET NULL,
    event_name VARCHAR(64) NOT NULL,
    properties JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analytics_event_name ON analytics_event(event_name);
CREATE INDEX IF NOT EXISTS idx_analytics_created_at ON analytics_event(created_at);
