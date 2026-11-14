package services

import (
	"conectainova/internal/domain/dto/avaliacao"
	"conectainova/internal/domain/factory"
	"conectainova/internal/domain/models"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"conectainova/internal/validator"
	"strings"
)

type IAvaliacaoService interface {
	CreateAvaliacao(*avaliacao.CreateAvaliacaoRequest) (*avaliacao.AvaliacaoResponse, bool, error)
	GetAvaliacoes(startupID, page, size uint) ([]*avaliacao.AvaliacaoResponse, error)
	UpdateAvaliacao(uint, *avaliacao.UpdateAvaliacaoRequest) (*avaliacao.AvaliacaoResponse, error)
}

type AvaliacaoService struct {
	AvaliacaoRepo   repositories.IAvaliacaoRepository
	NecessidadeRepo repositories.INecessidadeRepository
	SolucaoRepo     repositories.ISolucaoRepository
	Associacao      string
	Regiao          string
}

func NewAvaliacaoService(repo repositories.IAvaliacaoRepository, needs repositories.INecessidadeRepository, solutions repositories.ISolucaoRepository, associacao, regiao string) IAvaliacaoService {
	return &AvaliacaoService{AvaliacaoRepo: repo, NecessidadeRepo: needs, SolucaoRepo: solutions, Associacao: associacao, Regiao: regiao}
}

func (s *AvaliacaoService) CreateAvaliacao(request *avaliacao.CreateAvaliacaoRequest) (*avaliacao.AvaliacaoResponse, bool, error) {
	entity := factory.NewAvaliacao(request, s.Associacao, s.Regiao)
	if err := validator.ValidateAvaliacao(entity); err != nil {
		return nil, false, err
	}
	need, err := s.NecessidadeRepo.GetByID(entity.NecessidadeID)
	if err != nil {
		return nil, false, err
	}
	solution, err := s.SolucaoRepo.GetByID(entity.SolucaoID)
	if err != nil {
		return nil, false, err
	}
	if need.Categoria != solution.Categoria {
		return nil, false, models.NewErrorMessage("solucao", "A categoria da solução não atende ao problema.")
	}
	if strings.TrimSpace(need.ResultadoEsperado) == "" || len(need.CriteriosSucesso) == 0 {
		return nil, false, models.NewErrorMessage("briefing", "Complete o resultado esperado e os critérios de sucesso do problema antes de solicitar uma avaliação.")
	}
	saved, created, err := s.AvaliacaoRepo.CreateOrGet(entity)
	if err != nil {
		return nil, false, err
	}
	return mapper.AvaliacaoToResponse(saved), created, nil
}

func (s *AvaliacaoService) GetAvaliacoes(startupID, page, size uint) ([]*avaliacao.AvaliacaoResponse, error) {
	entities, err := s.AvaliacaoRepo.GetFiltered(startupID, page, size)
	if err != nil {
		return nil, err
	}
	return mapper.AvaliacoesToResponse(entities), nil
}

func (s *AvaliacaoService) UpdateAvaliacao(id uint, request *avaliacao.UpdateAvaliacaoRequest) (*avaliacao.AvaliacaoResponse, error) {
	entity, err := s.AvaliacaoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	updated := factory.UpdateAvaliacao(entity, request)
	if err := validator.ValidateAvaliacao(updated); err != nil {
		return nil, err
	}
	if err := s.AvaliacaoRepo.Update(updated); err != nil {
		return nil, err
	}
	return mapper.AvaliacaoToResponse(updated), nil
}
