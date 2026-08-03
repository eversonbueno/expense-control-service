package launches

import (
	"context"
	"expense-control-service/internal/entity"
	launchesRepo "expense-control-service/internal/repositories/launches"
)

type Launches interface {
	ListLaunches(ctx context.Context) ([]*entity.Lancamentos, error)
}

type launches struct {
	launchesRepo launchesRepo.Launches
}

func New(
	launchesRepo launchesRepo.Launches,
) Launches  {
	return &launches{launchesRepo: launchesRepo}
}

func (l *launches) ListLaunches(ctx context.Context) ([]*entity.Lancamentos, error)  {
	launches,err := l.launchesRepo.ListLaunches(ctx)
	if err != nil {
		return nil, err
	}

	return launches, nil
}

func (l launches) CreateLaunching(ctx context.Context, launching entity.Lancamentos) error  {
	newLAunching, err := l.launchesRepo.CreateLaunche(ctx, &launching)
	if err != nil{
		return err
	}

	if launching.FormaPagamento == "parcelado" && launching.ParcelasTotal > 1 {
		for i := 0; i <= launching.ParcelasTotal; i++  {
			valorParcela := launching.Valor/float64(launching.ParcelasTotal)

			parcelaLaunching := entity.Lancamentos{
				UsuarioID:             launching.UsuarioID,
				ContaID:               launching.ContaID,
				CategoriaID:           launching.CategoriaID,
				Tipo:                  "saida",
				Descricao:             "Manutenção Carro",
				Valor:                 valorParcela,
				FormaPagamento:        "parcelado",
				MetodoPagamento:       "cartão",
				DataCompra:            "2025-10-30",
				DataVencimento:        "2025-12-10",
				Status:                "pendente",
				ParcelasTotal:         launching.ParcelasTotal,
				ParcelaAtual:          i,
				IdParcelaPai:          newLAunching.IdParcelaPai,
			}
			if _, err := l.launchesRepo.CreateLaunche(ctx, &parcelaLaunching); err != nil {
				return err
			}
		}
	}

	return nil
}
