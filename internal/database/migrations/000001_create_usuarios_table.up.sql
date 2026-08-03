-- created_at/updated_at (em vez de criado_em) porque e o que
-- internal/repositories/users/queries.go e internal/entity/users.go
-- efetivamente consultam — ver reconciliacao de schema na 000007.
CREATE TABLE usuarios
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    nome       VARCHAR(100)        NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    senha_hash VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
