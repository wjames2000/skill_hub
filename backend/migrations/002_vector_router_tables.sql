-- Migration: Vector Embeddings and Router Logs
-- Description: Create tables for skill embeddings storage and router logging

BEGIN;

-- Skill embeddings table (PostgreSQL-side storage for vector data)
CREATE TABLE IF NOT EXISTS skill_embeddings (
    id BIGSERIAL PRIMARY KEY,
    skill_id BIGINT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    vector JSONB NOT NULL DEFAULT '[]',
    model_name VARCHAR(100) NOT NULL DEFAULT '',
    content_hash VARCHAR(64) NOT NULL DEFAULT '',
    chunk_index INT NOT NULL DEFAULT 0,
    chunk_text TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_skill_embeddings_skill_id ON skill_embeddings(skill_id);
CREATE INDEX IF NOT EXISTS idx_skill_embeddings_model_name ON skill_embeddings(model_name);
CREATE INDEX IF NOT EXISTS idx_skill_embeddings_content_hash ON skill_embeddings(content_hash);
CREATE UNIQUE INDEX IF NOT EXISTS idx_skill_embeddings_skill_chunk ON skill_embeddings(skill_id, chunk_index);

COMMENT ON TABLE skill_embeddings IS 'Stores skill embedding vectors and chunked text for vector search';
COMMENT ON COLUMN skill_embeddings.vector IS 'Float32 array as JSONB: the embedding vector';
COMMENT ON COLUMN skill_embeddings.model_name IS 'Name of the embedding model used';
COMMENT ON COLUMN skill_embeddings.content_hash IS 'Hash of source content for dedup';
COMMENT ON COLUMN skill_embeddings.chunk_index IS 'Order of chunk within skill';
COMMENT ON COLUMN skill_embeddings.chunk_text IS 'The text content of this chunk';

-- Router logs table (audit trail for routing decisions and executions)
CREATE TABLE IF NOT EXISTS router_logs (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL DEFAULT '',
    query TEXT NOT NULL DEFAULT '',
    query_embedding JSONB,
    matched_skill_id BIGINT NOT NULL DEFAULT 0,
    matched_skill_name VARCHAR(255) NOT NULL DEFAULT '',
    match_score DECIMAL(10,4) NOT NULL DEFAULT 0,
    match_strategy VARCHAR(50) NOT NULL DEFAULT '',
    is_executed BOOLEAN NOT NULL DEFAULT FALSE,
    execute_result TEXT,
    execute_duration INT NOT NULL DEFAULT 0,
    feedback_score SMALLINT NOT NULL DEFAULT 0,
    feedback_comment TEXT,
    user_id BIGINT NOT NULL DEFAULT 0,
    client_ip VARCHAR(45) NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_router_logs_session_id ON router_logs(session_id);
CREATE INDEX IF NOT EXISTS idx_router_logs_matched_skill_id ON router_logs(matched_skill_id);
CREATE INDEX IF NOT EXISTS idx_router_logs_user_id ON router_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_router_logs_created_at ON router_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_router_logs_feedback ON router_logs(feedback_score) WHERE feedback_score > 0;

COMMENT ON TABLE router_logs IS 'Audit trail for smart routing: match, execute, and feedback';
COMMENT ON COLUMN router_logs.session_id IS 'Session identifier for grouping related requests';
COMMENT ON COLUMN router_logs.query_embedding IS 'Optional: stored query vector for analysis';
COMMENT ON COLUMN router_logs.matched_skill_id IS 'ID of the matched skill';
COMMENT ON COLUMN router_logs.match_strategy IS 'Strategy used: vector, keyword, hybrid, hybrid+rerank';
COMMENT ON COLUMN router_logs.is_executed IS 'Whether the skill was actually executed';
COMMENT ON COLUMN router_logs.execute_result IS 'Result text from skill execution';
COMMENT ON COLUMN router_logs.execute_duration IS 'Execution duration in milliseconds';
COMMENT ON COLUMN router_logs.feedback_score IS 'User feedback: 1-5, 0 means no feedback';
COMMENT ON COLUMN router_logs.client_ip IS 'Client IP address for audit';

COMMIT;
