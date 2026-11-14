CREATE TABLE avaliacoes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    necessidade_id BIGINT NOT NULL REFERENCES necessidades(id) ON DELETE RESTRICT,
    solucao_id BIGINT NOT NULL REFERENCES solucoes(id) ON DELETE RESTRICT,
    mensagem TEXT NOT NULL DEFAULT '',
    resposta TEXT NOT NULL DEFAULT '',
    associacao VARCHAR(150) NOT NULL DEFAULT '',
    regiao VARCHAR(150) NOT NULL DEFAULT '',
    status VARCHAR(30) NOT NULL DEFAULT 'solicitada' CHECK (status IN ('solicitada', 'em_analise', 'piloto_proposto', 'nao_atendida', 'concluida')),
    CONSTRAINT idx_avaliacao_par UNIQUE (necessidade_id, solucao_id)
);
CREATE INDEX idx_avaliacoes_deleted_at ON avaliacoes (deleted_at);
CREATE INDEX idx_avaliacoes_solucao ON avaliacoes (solucao_id);
