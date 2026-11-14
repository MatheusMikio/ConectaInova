CREATE TABLE IF NOT EXISTS solucoes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    startup_id BIGINT NOT NULL REFERENCES startups(id) ON DELETE RESTRICT,
    nome VARCHAR(150) NOT NULL,
    descricao TEXT NOT NULL,
    categoria VARCHAR(20) NOT NULL CHECK (categoria IN ('estoque', 'vendas', 'atendimento', 'financeiro', 'marketing', 'logistica')),
    palavras_chave JSONB NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(palavras_chave) = 'array')
);

CREATE INDEX IF NOT EXISTS idx_solucoes_startup_id ON solucoes (startup_id);
CREATE INDEX IF NOT EXISTS idx_solucoes_categoria ON solucoes (categoria);
CREATE INDEX IF NOT EXISTS idx_solucoes_deleted_at ON solucoes (deleted_at);
