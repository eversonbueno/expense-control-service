CREATE TABLE orcamentos
(
    id           SERIAL PRIMARY KEY,
    usuario_id   INT            NOT NULL,
    categoria_id INT            NOT NULL,
    mes          INT            NOT NULL,
    ano          INT            NOT NULL,
    valor_limite DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (usuario_id) REFERENCES usuarios (id),
    FOREIGN KEY (categoria_id) REFERENCES categorias (id)
);