DROP TABLE IF EXISTS audit_log CASCADE;
DROP TABLE IF EXISTS risk_finding CASCADE;
DROP TABLE IF EXISTS document CASCADE;

DROP TRIGGER IF EXISTS trg_analysis_task_updated_at ON analysis_task;
DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS analysis_task CASCADE;
DROP TABLE IF EXISTS user_account CASCADE;
