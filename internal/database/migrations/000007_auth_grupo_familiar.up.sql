-- Migration: 000007_auth_grupo_familiar.up.sql
-- Autor: dev-dba
-- Projeto: expense-control-service
-- Demanda: Spec 001-autenticacao-grupo-familiar — cria grupo_familiar e convites,
--           adiciona email/senha_hash/grupo_familiar_id em usuarios.
--
-- ATENCAO: a tabela `usuarios` real neste banco NAO corresponde as migrations
-- 000001-000006 (nunca aplicadas — nao ha tabela schema_migrations). O schema
-- real e (id, nome, sobrenome, usuario, senha, saldo, created_at, updated_at).
-- Esta migration parte do schema REAL, nao do assumido pelas migrations antigas.

CREATE TABLE grupo_familiar
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE convites
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    codigo            VARCHAR(64)  NOT NULL,
    grupo_familiar_id INT          NOT NULL,
    criado_em         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expira_em         TIMESTAMP    NOT NULL,
    utilizado_em      TIMESTAMP    NULL,
    UNIQUE KEY uq_convites_codigo (codigo),
    CONSTRAINT fk_convites_grupo_familiar FOREIGN KEY (grupo_familiar_id) REFERENCES grupo_familiar (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Colunas novas em usuarios, nullable por enquanto para permitir o backfill
-- do usuario legado (id=1) antes de travar as constraints NOT NULL/UNIQUE.
ALTER TABLE usuarios
    ADD COLUMN email             VARCHAR(100) NULL AFTER usuario,
    ADD COLUMN senha_hash        VARCHAR(255) NULL AFTER senha,
    ADD COLUMN grupo_familiar_id INT          NULL AFTER saldo;

-- Backfill do usuario legado (id=1, Everson Bueno): cria um grupo familiar
-- proprio e popula email/senha_hash com uma senha temporaria (ver relatorio
-- do dev-dba para a senha em texto — trocar apos o primeiro login).
INSERT INTO grupo_familiar (criado_em) VALUES (CURRENT_TIMESTAMP);
SET @grupo_legado_id = LAST_INSERT_ID();

UPDATE usuarios
SET email             = 'eversonmbueno@gmail.com',
    senha_hash         = '$2y$10$NAbTAFfoNQnMYl4zIXoTMe7dHRbYWPHwhjtc8qU14ir0/3V5U6inu',
    grupo_familiar_id = @grupo_legado_id
WHERE id = 1;

-- Trava as constraints agora que todo registro existente tem os valores.
ALTER TABLE usuarios
    MODIFY COLUMN email             VARCHAR(100) NOT NULL,
    MODIFY COLUMN senha_hash        VARCHAR(255) NOT NULL,
    MODIFY COLUMN grupo_familiar_id INT          NOT NULL;

ALTER TABLE usuarios
    ADD CONSTRAINT uq_usuarios_email UNIQUE (email),
    ADD CONSTRAINT fk_usuarios_grupo_familiar FOREIGN KEY (grupo_familiar_id) REFERENCES grupo_familiar (id);
