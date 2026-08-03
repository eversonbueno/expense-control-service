-- Migration: 000007_auth_grupo_familiar.down.sql
-- Autor: dev-dba
-- Projeto: expense-control-service
-- Rollback da spec 001-autenticacao-grupo-familiar

ALTER TABLE usuarios
    DROP FOREIGN KEY fk_usuarios_grupo_familiar;

ALTER TABLE usuarios
    DROP COLUMN grupo_familiar_id;

DROP TABLE IF EXISTS convites;
DROP TABLE IF EXISTS grupo_familiar;
