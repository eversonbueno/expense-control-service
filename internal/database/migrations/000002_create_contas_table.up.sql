CREATE TABLE contas
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id    INT          NOT NULL,
    nome          VARCHAR(100) NOT NULL, -- Ex: "Nubank", "Carteira", "Bradesco"
    tipo          VARCHAR(50),           -- Ex: "corrente", "poupança", "cartão crédito"
    saldo_inicial DECIMAL(10, 2) DEFAULT 0,
    ativo         BOOLEAN        DEFAULT TRUE,
    FOREIGN KEY (usuario_id) REFERENCES usuarios (id)
);