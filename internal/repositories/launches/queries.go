package launches

var(
	ListLaunches = `
		SELECT 
			id, 
			idfk_usuario AS usuario, 
			idfk_forma_pagamento AS forma_pagamento, 
			idfk_tipo_lancamento AS tipo_lancamento, 
			idfk_categoria_lancamento AS categoria_lancamento, 
			mes, 
			ano, 
			parcelado, 
			parcelado_quantidade, 
			descricao, 
			valor, 
			created_at, 
			updated_at
		FROM tipo_lancamento
	`
)