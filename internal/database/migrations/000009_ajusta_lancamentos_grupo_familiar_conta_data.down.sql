-- Migration: 000009_ajusta_lancamentos_grupo_familiar_conta_data.down.sql
-- Rollback dos ajustes de grupo_familiar_id/conta_id/data/excluido em `lancamentos`.

ALTER TABLE lancamentos
    DROP FOREIGN KEY fk_lancamentos_grupo_familiar,
    DROP FOREIGN KEY fk_lancamentos_conta,
    DROP COLUMN grupo_familiar_id,
    DROP COLUMN conta_id,
    DROP COLUMN data,
    DROP COLUMN excluido;
