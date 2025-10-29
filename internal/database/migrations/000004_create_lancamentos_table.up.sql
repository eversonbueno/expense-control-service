CREATE TABLE lancamentos
(
    id                     INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id             INT            NOT NULL,
    conta_id               INT            NOT NULL,
    categoria_id           INT,
    tipo                   VARCHAR(20)    NOT NULL,       -- "entrada" ou "saida"
    descricao              TEXT,
    valor                  DECIMAL(10, 2) NOT NULL,
    forma_pagamento        VARCHAR(50),                   -- "avista", "parcelado"
    metodo_pagamento       VARCHAR(50),                   -- "pix", "cartao", "dinheiro", etc.
    data_compra            DATE           NOT NULL,
    data_vencimento        DATE,
    data_pagamento         DATE,
    status                 VARCHAR(20) DEFAULT 'pendente',-- "pendente", "pago", "cancelado"
    recorrente             BOOLEAN     DEFAULT FALSE,
    frequencia_recorrencia VARCHAR(20),                   -- "mensal", "anual", etc.
    parcelas_total         INT,
    parcela_atual          INT,
    id_parcela_pai         INT,                           -- ligação com o lançamento pai
    comprovante_url        TEXT,
    criado_em              TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    atualizado_em          TIMESTAMP   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios (id),
    FOREIGN KEY (conta_id) REFERENCES contas (id),
    FOREIGN KEY (categoria_id) REFERENCES categorias (id),
    FOREIGN KEY (id_parcela_pai) REFERENCES lancamentos (id)
);