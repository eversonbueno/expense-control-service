-- Migration: 000010_ajusta_categoria_lancamento_tipo.up.sql
-- Autor: dev-dba
-- Projeto: expense-control-service
-- Demanda: Spec BE-003 (Lançamentos Financeiros), RN-06 — categoria deve ser compatível com o tipo
--   do lançamento (ex: categoria de despesa não pode ser usada num lançamento do tipo receita).
--   `categoria_lancamento` hoje não tem nenhum vínculo com `tipo_lancamento` — adiciona esse vínculo.
-- Data: 2026-08-03

ALTER TABLE categoria_lancamento
    ADD COLUMN idfk_tipo_lancamento INT NULL AFTER descricao;

-- Backfill: as 6 categorias existentes (Alimentação, Lazer, Combustivel, Saude, Carro, Outros) são,
-- pelo nome, todas categorias de despesa — atribuídas a "Saida" (id=2). Não há como inferir do dado
-- existente se alguma delas deveria ser categoria de receita; ajustar manualmente se necessário.
UPDATE categoria_lancamento SET idfk_tipo_lancamento = 2;

ALTER TABLE categoria_lancamento
    MODIFY COLUMN idfk_tipo_lancamento INT NOT NULL,
    ADD CONSTRAINT fk_categoria_lancamento_tipo FOREIGN KEY (idfk_tipo_lancamento) REFERENCES tipo_lancamento (id);
