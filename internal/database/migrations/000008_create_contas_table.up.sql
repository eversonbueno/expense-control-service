-- Migration: 000008_create_contas_table.up.sql
-- Autor: dev-dba
-- Projeto: expense-control-service
-- Demanda: Spec BE-002 (Cadastro de Contas) — a tabela `contas` nunca existiu de fato neste banco
--   (os arquivos de migration 000002_create_contas_table nunca foram aplicados aqui; o schema real
--   diverge deles). Criação do zero, já no modelo correto (grupo_familiar_id, não usuario_id).
-- Data: 2026-08-03

CREATE TABLE contas
(
    id                INT            NOT NULL AUTO_INCREMENT,
    grupo_familiar_id INT            NOT NULL,
    nome              VARCHAR(100)   NOT NULL,
    tipo              ENUM ('corrente', 'poupanca', 'cartao_credito', 'dinheiro') NOT NULL,
    fechamento_cartao TINYINT UNSIGNED NULL,
    saldo_inicial     DECIMAL(10, 2) NOT NULL DEFAULT 0,
    ativo             BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY fk_contas_grupo_familiar (grupo_familiar_id),
    CONSTRAINT fk_contas_grupo_familiar FOREIGN KEY (grupo_familiar_id) REFERENCES grupo_familiar (id),
    CONSTRAINT chk_contas_fechamento_cartao
        CHECK (fechamento_cartao IS NULL OR (tipo = 'cartao_credito' AND fechamento_cartao BETWEEN 1 AND 31))
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Conta placeholder para o único lançamento legado (id=1, grupo_familiar_id=1), que vai precisar de
-- uma conta_id válida na migration 000009 (RN-04 da spec 003 exige conta obrigatória e ativa).
INSERT INTO contas (grupo_familiar_id, nome, tipo, saldo_inicial, ativo)
VALUES (1, 'Conta Legada (migração)', 'dinheiro', 0, TRUE);
