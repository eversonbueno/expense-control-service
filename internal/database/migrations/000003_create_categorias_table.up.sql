CREATE TABLE categorias
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id INT          NOT NULL,
    nome       VARCHAR(100) NOT NULL, -- Ex: "Alimentação"
    tipo       VARCHAR(20)  NOT NULL, -- "entrada" ou "saida"
    cor        VARCHAR(10),           -- Ex: "#FF5733" (para gráficos)
    FOREIGN KEY (usuario_id) REFERENCES usuarios (id)
);