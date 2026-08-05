-- Migration: 000009_ajusta_lancamentos_grupo_familiar_conta_data.up.sql
-- Autor: dev-dba
-- Projeto: expense-control-service
-- Demanda: Spec BE-003 (Lançamentos Financeiros) — a tabela `lancamentos` real (idfk_usuario,
--   idfk_forma_pagamento, idfk_tipo_lancamento, idfk_categoria_lancamento, mes, ano, parcelado,
--   parcelado_quantidade, descricao, valor) já bate com o código Go existente, mas precisa de:
--   grupo_familiar_id (isolamento por família, RN-01), conta_id (RN-04, depende da 000008),
--   data completa em vez de mes/ano (RN-09) e soft-delete (RN-11).
-- Data: 2026-08-03

ALTER TABLE lancamentos
    ADD COLUMN grupo_familiar_id INT NULL AFTER idfk_usuario,
    ADD COLUMN conta_id INT NULL AFTER grupo_familiar_id,
    ADD COLUMN data DATE NULL AFTER ano,
    ADD COLUMN excluido BOOLEAN NOT NULL DEFAULT FALSE AFTER valor;

-- Backfill grupo_familiar_id a partir do usuário dono do lançamento (mesmo padrão da migration 000007)
UPDATE lancamentos l
    JOIN usuarios u ON u.id = l.idfk_usuario
    SET l.grupo_familiar_id = u.grupo_familiar_id;

-- Backfill conta_id para a "Conta Legada" criada na migration 000008
UPDATE lancamentos l
    JOIN contas c ON c.grupo_familiar_id = l.grupo_familiar_id AND c.nome = 'Conta Legada (migração)'
    SET l.conta_id = c.id;

-- Backfill data a partir de mes/ano — dia 01 como aproximação (dia exato do lançamento legado é
-- desconhecido; decisão confirmada com o usuário nesta rodada de migration).
UPDATE lancamentos
    SET data = STR_TO_DATE(CONCAT(ano, '-', LPAD(mes, 2, '0'), '-01'), '%Y-%m-%d')
    WHERE mes IS NOT NULL AND ano IS NOT NULL;

ALTER TABLE lancamentos
    MODIFY COLUMN grupo_familiar_id INT NOT NULL,
    MODIFY COLUMN conta_id INT NOT NULL,
    MODIFY COLUMN data DATE NOT NULL,
    ADD CONSTRAINT fk_lancamentos_grupo_familiar FOREIGN KEY (grupo_familiar_id) REFERENCES grupo_familiar (id),
    ADD CONSTRAINT fk_lancamentos_conta FOREIGN KEY (conta_id) REFERENCES contas (id);
