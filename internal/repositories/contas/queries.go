package contas

var (
	CreateConta = `
		INSERT INTO contas (grupo_familiar_id, nome, tipo, fechamento_cartao, saldo_inicial)
		VALUES (?, ?, ?, ?, ?)
	`

	FindContaById = `
		SELECT
			id,
			grupo_familiar_id,
			nome,
			tipo,
			fechamento_cartao,
			saldo_inicial,
			ativo,
			created_at,
			updated_at
		FROM contas
		WHERE id = ?
	`

	ListContasByGrupoFamiliar = `
		SELECT
			id,
			grupo_familiar_id,
			nome,
			tipo,
			fechamento_cartao,
			saldo_inicial,
			ativo,
			created_at,
			updated_at
		FROM contas
		WHERE grupo_familiar_id = ?
		ORDER BY nome
	`

	UpdateConta = `
		UPDATE contas
		SET nome = ?, tipo = ?, fechamento_cartao = ?, saldo_inicial = ?
		WHERE id = ?
	`

	SetContaAtivo = `
		UPDATE contas
		SET ativo = ?
		WHERE id = ?
	`
)
