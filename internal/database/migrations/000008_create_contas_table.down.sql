-- Migration: 000008_create_contas_table.down.sql
-- Rollback da criação da tabela `contas` (spec BE-002).
-- Só pode ser executado depois do rollback da 000009 (que remove a FK lancamentos.conta_id -> contas).

DROP TABLE IF EXISTS contas;
