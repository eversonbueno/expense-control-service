package users

var(
	ListUserById = `
		SELECT 
			nome, 
			sobrenome
		FROM usuarios
		WHERE id = ?
	`
)
