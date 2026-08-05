package launches

var (
	launchColumns = `
		id,
		idfk_usuario,
		grupo_familiar_id,
		conta_id,
		idfk_forma_pagamento,
		idfk_tipo_lancamento,
		idfk_categoria_lancamento,
		mes,
		ano,
		data,
		parcelado,
		parcelado_quantidade,
		descricao,
		valor,
		excluido,
		created_at,
		updated_at
	`

	ListLaunches = `
		SELECT ` + launchColumns + `
		FROM lancamentos
		WHERE grupo_familiar_id = ?
			AND excluido = FALSE
			AND (? IS NULL OR MONTH(data) = ?)
			AND (? IS NULL OR YEAR(data) = ?)
		ORDER BY data DESC, id DESC
	`

	FindLaunchById = `
		SELECT ` + launchColumns + `
		FROM lancamentos
		WHERE id = ?
	`

	CreateLaunch = `
		INSERT INTO lancamentos
			(idfk_usuario, grupo_familiar_id, conta_id, idfk_forma_pagamento, idfk_tipo_lancamento, idfk_categoria_lancamento, data, descricao, valor, parcelado, parcelado_quantidade)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0)
	`

	UpdateLaunch = `
		UPDATE lancamentos
		SET conta_id = ?, idfk_forma_pagamento = ?, idfk_tipo_lancamento = ?, idfk_categoria_lancamento = ?, data = ?, descricao = ?, valor = ?
		WHERE id = ?
	`

	SoftDeleteLaunch = `
		UPDATE lancamentos
		SET excluido = TRUE
		WHERE id = ?
	`
)
