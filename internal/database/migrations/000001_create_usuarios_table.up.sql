CREATE TABLE usuarios
(
    id         SERIAL PRIMARY KEY,
    nome       VARCHAR(100)        NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    senha_hash VARCHAR(255)        NOT NULL,
    criado_em  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
