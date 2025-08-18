package users

var(
	ListUsers = `
		SELECT 
			id, 
			nome, 
			sobrenome, 
			usuario, 
			senha,
			saldo,
			created_at, 
			updated_at
		FROM users
	`
)
