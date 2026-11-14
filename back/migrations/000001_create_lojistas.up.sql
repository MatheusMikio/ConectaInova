CREATE TABLE IF NOT EXISTS lojistas (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    nome VARCHAR(150) NOT NULL,
    segmento VARCHAR(150) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lojistas_deleted_at ON lojistas (deleted_at);
