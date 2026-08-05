package launch_category

var (
	ListLaunchCategory = `
		SELECT
			id,
			descricao
		FROM categoria_lancamento
	`

	FindLaunchCategoryById = `
		SELECT
			id,
			descricao,
			idfk_tipo_lancamento
		FROM categoria_lancamento
		WHERE id = ?
	`
)
