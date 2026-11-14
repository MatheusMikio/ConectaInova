CREATE TABLE IF NOT EXISTS matches (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    necessidade_id BIGINT NOT NULL REFERENCES necessidades(id) ON DELETE RESTRICT,
    solucao_id BIGINT NOT NULL REFERENCES solucoes(id) ON DELETE RESTRICT,
    pontuacao NUMERIC NOT NULL CHECK (pontuacao BETWEEN 0 AND 100),
    CONSTRAINT idx_match_pair UNIQUE (necessidade_id, solucao_id)
);

CREATE INDEX IF NOT EXISTS idx_matches_deleted_at ON matches (deleted_at);
