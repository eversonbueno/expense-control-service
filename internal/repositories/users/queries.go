package users

var (
	ListUserById = `
		SELECT
			id,
			nome,
			email,
			senha_hash,
			grupo_familiar_id,
			created_at,
			updated_at
		FROM usuarios
		WHERE id = ?
	`

	FindByEmail = `
		SELECT
			id,
			nome,
			email,
			senha_hash,
			grupo_familiar_id,
			created_at,
			updated_at
		FROM usuarios
		WHERE email = ?
	`

	CreateUser = `
		INSERT INTO usuarios (nome, email, senha_hash, grupo_familiar_id)
		VALUES (?, ?, ?, ?)
	`

	ListUsers = `
		SELECT
			id,
			nome,
			email,
			senha_hash,
			grupo_familiar_id,
			created_at,
			updated_at
		FROM usuarios
	`
)
