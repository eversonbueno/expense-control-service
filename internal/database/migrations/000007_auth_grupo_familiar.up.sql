-- Migration: 000007_auth_grupo_familiar.up.sql
-- Autor: dev-dba
-- Projeto: expense-control-service
-- Demanda: Spec 001-autenticacao-grupo-familiar — cria grupo_familiar e convites,
--           adiciona grupo_familiar_id em usuarios.
--
-- Reconciliacao de schema (2026-08-03): esta migration assumia um schema
-- legado de `usuarios` (sobrenome, usuario, senha, saldo) nunca confirmado e
-- incompativel com a cadeia 000001-000005 deste repositorio nem com o que
-- internal/repositories/users/queries.go realmente consulta. A 000001 foi
-- corrigida para criar diretamente (id, nome, email, senha_hash, created_at,
-- updated_at), entao aqui so falta adicionar grupo_familiar_id.

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
    expira_em         DATETIME     NOT NULL,
    utilizado_em      TIMESTAMP    NULL,
    UNIQUE KEY uq_convites_codigo (codigo),
    CONSTRAINT fk_convites_grupo_familiar FOREIGN KEY (grupo_familiar_id) REFERENCES grupo_familiar (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Coluna nova em usuarios, nullable por enquanto para permitir o backfill
-- do usuario legado (id=1) antes de travar a constraint NOT NULL.
ALTER TABLE usuarios
    ADD COLUMN grupo_familiar_id INT NULL;

-- Backfill do usuario legado (id=1, Everson Bueno), se existir: cria um
-- grupo familiar proprio e associa. Em bancos novos, sem esse usuario,
-- as duas instrucoes abaixo nao afetam nenhuma linha (no-op seguro).
INSERT INTO grupo_familiar (criado_em)
SELECT CURRENT_TIMESTAMP FROM usuarios WHERE id = 1;

UPDATE usuarios
SET grupo_familiar_id = LAST_INSERT_ID()
WHERE id = 1;

-- Trava a constraint agora que todo registro existente tem o valor.
ALTER TABLE usuarios
    MODIFY COLUMN grupo_familiar_id INT NOT NULL;

ALTER TABLE usuarios
    ADD CONSTRAINT fk_usuarios_grupo_familiar FOREIGN KEY (grupo_familiar_id) REFERENCES grupo_familiar (id);
