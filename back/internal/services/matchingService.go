package services

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	"conectainova/internal/domain/dto/matching"
	"conectainova/internal/domain/schema"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type MatchingService struct {
	NecessidadeRepo repositories.INecessidadeRepository
	SolucaoRepo     repositories.ISolucaoRepository
}

type IMatchingService interface {
	GetMatches(id uint) (*matching.MatchingResponse, error)
}

func NewMatchingService(necessidadeRepo repositories.INecessidadeRepository, solucaoRepo repositories.ISolucaoRepository) IMatchingService {
	return &MatchingService{NecessidadeRepo: necessidadeRepo, SolucaoRepo: solucaoRepo}
}

func normalize(value string) string {
	result, _, _ := transform.String(transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC), strings.ToLower(value))
	return strings.Join(strings.FieldsFunc(result, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}), " ")
}

// Same-category solutions start at 70. Keyword coverage contributes up to 30.
func Score(necessidade schema.Necessidade, solucao schema.Solucao) float64 {
	if necessidade.Categoria != solucao.Categoria {
		return 0
	}
	keywords := map[string]bool{}
	for _, keyword := range solucao.PalavrasChave {
		if normalized := normalize(keyword); normalized != "" {
			keywords[normalized] = true
		}
	}
	if len(keywords) == 0 {
		return 70
	}
	content := " " + normalize(necessidade.Titulo+" "+necessidade.Descricao) + " "
	hits := 0
	for keyword := range keywords {
		if strings.Contains(content, " "+keyword+" ") {
			hits++
		}
	}
	return math.Round((70+30*float64(hits)/float64(len(keywords)))*100) / 100
}

func explainMatch(need schema.Necessidade, solution schema.Solucao) ([]string, []string) {
	content := " " + normalize(need.Titulo+" "+need.Descricao) + " "
	found := []string{}
	seen := map[string]bool{}
	for _, keyword := range solution.PalavrasChave {
		word := normalize(keyword)
		if word != "" && !seen[word] && strings.Contains(content, " "+word+" ") {
			found = append(found, keyword)
			seen[word] = true
		}
	}
	reasons := []string{fmt.Sprintf("A solução atende à mesma categoria do problema: %s.", need.Categoria)}
	if len(found) > 0 {
		reasons = append(reasons, "Termos da sua necessidade aparecem nas palavras-chave da solução: "+strings.Join(found, ", ")+".")
	} else {
		reasons = append(reasons, "Não foram encontradas palavras-chave em comum; a recomendação considera apenas a categoria.")
	}
	return reasons, found
}

func (s *MatchingService) GetMatches(id uint) (*matching.MatchingResponse, error) {
	necessidade, err := s.NecessidadeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	solucoes, err := s.SolucaoRepo.GetByCategoria(necessidade.Categoria)
	if err != nil {
		return nil, err
	}
	result := &matching.MatchingResponse{NecessidadeID: id, Matches: []matching.MatchResponse{}}
	for _, solucao := range solucoes {
		if solucao.Startup.ID == 0 {
			continue
		}
		response := mapper.MatchToResponse(solucao, Score(*necessidade, *solucao))
		response.Motivos, response.PalavrasEncontradas = explainMatch(*necessidade, *solucao)
		response.PontosAValidar = []string{"Preço e condições comerciais.", "Prazo de implantação e suporte.", "Integrações com as ferramentas da loja.", "Atendimento aos critérios de sucesso definidos pelo lojista."}
		result.Matches = append(result.Matches, *response)
	}
	sort.Slice(result.Matches, func(i, j int) bool {
		if result.Matches[i].Pontuacao == result.Matches[j].Pontuacao {
			return result.Matches[i].Solucao.ID < result.Matches[j].Solucao.ID
		}
		return result.Matches[i].Pontuacao > result.Matches[j].Pontuacao
	})
	result.SugerirDesafio = len(result.Matches) == 0
	return result, nil
}
