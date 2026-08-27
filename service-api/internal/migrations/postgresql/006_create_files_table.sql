CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id BIGINT REFERENCES users(tg_id),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    version BIGINT NOT NULL DEFAULT 0,

    local_path VARCHAR NOT NULL,

    tg_file_id VARCHAR(128),
    tg_file_unique_id VARCHAR(64) UNIQUE,
    file_type VARCHAR(20),

    mime_type VARCHAR(100),
    file_size BIGINT,
    width INT,
    height INT,
    duration INT,
    
);