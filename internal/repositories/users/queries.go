package users

var(
	ListUsers = `
		SELECT 
			id, 
			nome_completo, 
			cpf_cnpj, 
			email, 
			password, 
			tipo_usuario, 
			saldo_usuario
		FROM usuarios
	`
)
