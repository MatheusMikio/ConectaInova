CREATE TABLE IF NOT EXISTS startups (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    nome VARCHAR(150) NOT NULL,
    descricao TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_startups_deleted_at ON startups (deleted_at);
