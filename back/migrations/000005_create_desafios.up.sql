CREATE TABLE IF NOT EXISTS desafios (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    necessidade_id BIGINT NOT NULL REFERENCES necessidades(id) ON DELETE RESTRICT,
    titulo VARCHAR(200) NOT NULL,
    descricao TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'aberto' CHECK (status IN ('aberto', 'em_andamento', 'encerrado'))
);

CREATE INDEX IF NOT EXISTS idx_desafios_necessidade_id ON desafios (necessidade_id);
CREATE INDEX IF NOT EXISTS idx_desafios_deleted_at ON desafios (deleted_at);
