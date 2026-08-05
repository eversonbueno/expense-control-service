package entity

type LaunchCategory struct {
	ID               uint   `json:"id"`
	Descricao        string `json:"descricao"`
	TipoLancamentoID uint   `json:"idfk_tipo_lancamento"`
}
