package entity

type LauchType struct {
	ID        uint   `json:"id"`
	Descricao string `json:"descricao"`
}

var (
	LauchTypeEntrada = LauchType{ID: 1, Descricao: "ENTRADA"}
	LauchTypeSaida   = LauchType{ID: 2, Descricao: "SAIDA"}
)
