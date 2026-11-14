CREATE TABLE IF NOT EXISTS necessidades (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    lojista_id BIGINT NOT NULL REFERENCES lojistas(id) ON DELETE RESTRICT,
    titulo VARCHAR(200) NOT NULL,
    descricao TEXT NOT NULL,
    categoria VARCHAR(20) NOT NULL CHECK (categoria IN ('estoque', 'vendas', 'atendimento', 'financeiro', 'marketing', 'logistica')),
    prioridade VARCHAR(10) NOT NULL CHECK (prioridade IN ('baixa', 'media', 'alta')),
    status VARCHAR(20) NOT NULL DEFAULT 'aberta' CHECK (status IN ('aberta', 'em_analise', 'resolvida'))
);

CREATE INDEX IF NOT EXISTS idx_necessidades_lojista_id ON necessidades (lojista_id);
CREATE INDEX IF NOT EXISTS idx_necessidades_deleted_at ON necessidades (deleted_at);
