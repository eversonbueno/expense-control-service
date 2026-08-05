-- Migration: 000010_ajusta_categoria_lancamento_tipo.down.sql
-- Rollback do vínculo categoria_lancamento -> tipo_lancamento.

ALTER TABLE categoria_lancamento
    DROP FOREIGN KEY fk_categoria_lancamento_tipo,
    DROP COLUMN idfk_tipo_lancamento;
