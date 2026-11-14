package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

type publicProblemServiceStub struct {
	services.INecessidadeService
	page uint
	size uint
}

func (s *publicProblemServiceStub) GetNecessidades(page, size uint) ([]*necessidade.NecessidadeListResponse, error) {
	s.page, s.size = page, size
	return []*necessidade.NecessidadeListResponse{{LojistaNome: "Loja"}}, nil
}

func TestListProblemsWithoutIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &publicProblemServiceStub{}
	router := gin.New()
	router.GET("/api/v1/necessidades", ListNecessidadesHandler(service))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest("GET", "/api/v1/necessidades?page=2&size=12", nil))
	if response.Code != 200 || service.page != 2 || service.size != 12 {
		t.Fatalf("unexpected listing response: %d %s", response.Code, response.Body)
	}
	var payload struct {
		Data []necessidade.NecessidadeListResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Data) != 1 || payload.Data[0].LojistaNome != "Loja" {
		t.Fatal("missing store in response")
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest("GET", "/api/v1/necessidades?page=invalid", nil))
	if response.Code != 400 {
		t.Fatal("expected invalid pagination to be rejected")
	}
}
