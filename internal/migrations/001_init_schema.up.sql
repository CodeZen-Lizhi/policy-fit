CREATE TABLE IF NOT EXISTS user_account (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(32) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analysis_task (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_account(id),
    request_id VARCHAR(64),
    status VARCHAR(32) NOT NULL,
    risk_summary JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE analysis_task
    ADD COLUMN IF NOT EXISTS request_id VARCHAR(64);

ALTER TABLE analysis_task
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();

CREATE UNIQUE INDEX IF NOT EXISTS idx_task_request_id_unique
    ON analysis_task(request_id)
    WHERE request_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_task_user_id ON analysis_task(user_id);
CREATE INDEX IF NOT EXISTS idx_task_status ON analysis_task(status);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_analysis_task_updated_at ON analysis_task;

CREATE TRIGGER trg_analysis_task_updated_at
    BEFORE UPDATE ON analysis_task
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS document (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES analysis_task(id) ON DELETE CASCADE,
    doc_type VARCHAR(32) NOT NULL,
    file_name VARCHAR(256) NOT NULL,
    storage_key VARCHAR(512) NOT NULL,
    parse_status VARCHAR(32) NOT NULL,
    parsed_text TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_document_doc_type'
    ) THEN
        ALTER TABLE document
            ADD CONSTRAINT chk_document_doc_type
            CHECK (doc_type IN ('report', 'policy', 'disclosure'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_document_task_id ON document(task_id);
CREATE INDEX IF NOT EXISTS idx_document_parse_status ON document(parse_status);

CREATE TABLE IF NOT EXISTS risk_finding (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES analysis_task(id) ON DELETE CASCADE,
    level VARCHAR(16) NOT NULL,
    topic VARCHAR(64) NOT NULL,
    summary TEXT NOT NULL,
    health_evidence JSONB NOT NULL,
    policy_evidence JSONB NOT NULL,
    questions JSONB NOT NULL,
    actions JSONB,
    confidence NUMERIC(4,3),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_risk_finding_level'
    ) THEN
        ALTER TABLE risk_finding
            ADD CONSTRAINT chk_risk_finding_level
            CHECK (level IN ('red', 'yellow', 'green'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_finding_task_id ON risk_finding(task_id);
CREATE INDEX IF NOT EXISTS idx_finding_level ON risk_finding(level);
CREATE INDEX IF NOT EXISTS idx_finding_topic ON risk_finding(topic);

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT REFERENCES analysis_task(id) ON DELETE SET NULL,
    actor_id BIGINT REFERENCES user_account(id) ON DELETE SET NULL,
    action VARCHAR(64) NOT NULL,
    target_type VARCHAR(64) NOT NULL,
    target_id VARCHAR(128),
    detail JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_task_id ON audit_log(task_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_action ON audit_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log(created_at);
