package grupo_familiar

var (
	CreateGrupoFamiliar = `
		INSERT INTO grupo_familiar (criado_em)
		VALUES (NOW())
	`

	GetGrupoFamiliarById = `
		SELECT id, criado_em
		FROM grupo_familiar
		WHERE id = ?
	`
)
