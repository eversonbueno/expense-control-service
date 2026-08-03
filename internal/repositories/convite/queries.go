package convite

var (
	CreateConvite = `
		INSERT INTO convites (codigo, grupo_familiar_id, expira_em)
		VALUES (?, ?, ?)
	`

	FindByCodigo = `
		SELECT id, codigo, grupo_familiar_id, criado_em, expira_em, utilizado_em
		FROM convites
		WHERE codigo = ?
	`

	MarcarUtilizado = `
		UPDATE convites
		SET utilizado_em = ?
		WHERE id = ? AND utilizado_em IS NULL
	`
)
