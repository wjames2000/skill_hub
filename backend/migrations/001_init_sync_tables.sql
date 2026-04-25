-- Migration: 001_init_sync_tables
-- Description: Create initial tables for skill sync module

CREATE TABLE IF NOT EXISTS skills (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    display_name    VARCHAR(255),
    description     TEXT,
    version         VARCHAR(50),
    author          VARCHAR(255),
    repository      VARCHAR(512) NOT NULL,
    repo_owner      VARCHAR(255),
    repo_name       VARCHAR(255),
    default_branch  VARCHAR(255),
    skill_path      VARCHAR(512),
    skill_file_sha  VARCHAR(64),
    avatar_url      VARCHAR(512),
    homepage        VARCHAR(512),
    license         VARCHAR(100),
    stars           INT DEFAULT 0,
    forks           INT DEFAULT 0,
    open_issues     INT DEFAULT 0,
    language        VARCHAR(100),
    topics          JSONB DEFAULT '[]',
    category        VARCHAR(100),
    tags            JSONB DEFAULT '[]',
    readme          TEXT,
    installs        BIGINT DEFAULT 0,
    score           DECIMAL(10,2) DEFAULT 0,
    is_official     BOOLEAN DEFAULT FALSE,
    is_archived     BOOLEAN DEFAULT FALSE,
    scan_passed     BOOLEAN DEFAULT TRUE,
    scan_report     TEXT,
    status          SMALLINT DEFAULT 1,
    extra           JSONB DEFAULT '{}',
    last_sync_at    TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_skills_repository ON skills(repository);
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);
CREATE INDEX IF NOT EXISTS idx_skills_category ON skills(category);
CREATE INDEX IF NOT EXISTS idx_skills_stars ON skills(stars DESC);
CREATE INDEX IF NOT EXISTS idx_skills_status ON skills(status);
CREATE INDEX IF NOT EXISTS idx_skills_scan_passed ON skills(scan_passed);
CREATE INDEX IF NOT EXISTS idx_skills_last_sync_at ON skills(last_sync_at);

CREATE TABLE IF NOT EXISTS skill_versions (
    id          BIGSERIAL PRIMARY KEY,
    skill_id    INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    version     VARCHAR(50),
    skill_sha   VARCHAR(64),
    change_log  TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_skill_versions_skill_id ON skill_versions(skill_id);

CREATE TABLE IF NOT EXISTS skill_categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon        VARCHAR(255),
    parent_id   INT DEFAULT 0,
    sort_order  INT DEFAULT 0,
    skill_count INT DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sync_tasks (
    id              BIGSERIAL PRIMARY KEY,
    type            VARCHAR(20) NOT NULL,
    strategy        VARCHAR(30),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_repos     INT DEFAULT 0,
    found_repos     INT DEFAULT 0,
    parsed_skills   INT DEFAULT 0,
    new_skills      INT DEFAULT 0,
    updated_skills  INT DEFAULT 0,
    failed_skills   INT DEFAULT 0,
    scanned_repos   INT DEFAULT 0,
    error_count     INT DEFAULT 0,
    error_message   TEXT,
    started_at      TIMESTAMP WITH TIME ZONE,
    finished_at     TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sync_tasks_type ON sync_tasks(type);
CREATE INDEX IF NOT EXISTS idx_sync_tasks_status ON sync_tasks(status);
CREATE INDEX IF NOT EXISTS idx_sync_tasks_created_at ON sync_tasks(created_at DESC);
