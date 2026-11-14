package container

import (
	"conectainova/config"
	"conectainova/internal/repositories"
	"conectainova/internal/services"
	"gorm.io/gorm"
)

type Container struct {
	LojistaService     services.ILojistaService
	NecessidadeService services.INecessidadeService
	StartupService     services.IStartupService
	SolucaoService     services.ISolucaoService
	DesafioService     services.IDesafioService
	MatchingService    services.IMatchingService
	AvaliacaoService   services.IAvaliacaoService
}

func NewContainer(db *gorm.DB) *Container {
	lojistaRepo := repositories.NewLojistaRepository(db)
	necessidadeRepo := repositories.NewNecessidadeRepository(db)
	startupRepo := repositories.NewStartupRepository(db)
	solucaoRepo := repositories.NewSolucaoRepository(db)
	desafioRepo := repositories.NewDesafioRepository(db)
	avaliacaoRepo := repositories.NewAvaliacaoRepository(db)

	return &Container{
		LojistaService:     services.NewLojistaService(lojistaRepo),
		NecessidadeService: services.NewNecessidadeService(necessidadeRepo, lojistaRepo),
		StartupService:     services.NewStartupService(startupRepo),
		SolucaoService:     services.NewSolucaoService(solucaoRepo, startupRepo),
		DesafioService:     services.NewDesafioService(desafioRepo, necessidadeRepo),
		MatchingService:    services.NewMatchingService(necessidadeRepo, solucaoRepo),
		AvaliacaoService:   services.NewAvaliacaoService(avaliacaoRepo, necessidadeRepo, solucaoRepo, config.GetAssociacaoNome(), config.GetAssociacaoRegiao()),
	}
}
