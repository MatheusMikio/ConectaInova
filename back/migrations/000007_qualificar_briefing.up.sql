ALTER TABLE necessidades
    ADD COLUMN impacto TEXT NOT NULL DEFAULT '',
    ADD COLUMN resultado_esperado TEXT NOT NULL DEFAULT '',
    ADD COLUMN restricoes TEXT NOT NULL DEFAULT '',
    ADD COLUMN canal_venda VARCHAR(10) NOT NULL DEFAULT '',
    ADD COLUMN criterios_sucesso JSONB NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(criterios_sucesso) = 'array');

ALTER TABLE desafios
    ADD COLUMN criterios_sucesso JSONB NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(criterios_sucesso) = 'array'),
    ADD COLUMN prazo_dias INTEGER NOT NULL DEFAULT 7 CHECK (prazo_dias BETWEEN 1 AND 365);
