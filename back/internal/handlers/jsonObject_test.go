package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMutationHandlersRejectNullBeforeCallingServices(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handlers := []struct {
		name   string
		create gin.HandlerFunc
		update gin.HandlerFunc
	}{
		{"lojista", CreateLojistaHandler(nil), UpdateLojistaHandler(nil)},
		{"startup", CreateStartupHandler(nil), UpdateStartupHandler(nil)},
		{"solucao", CreateSolucaoHandler(nil), UpdateSolucaoHandler(nil)},
		{"necessidade", CreateNecessidadeHandler(nil), UpdateNecessidadeHandler(nil)},
		{"desafio", CreateDesafioHandler(nil), UpdateDesafioHandler(nil)},
		{"avaliacao", CreateAvaliacaoHandler(nil), UpdateAvaliacaoHandler(nil)},
	}
	for _, tt := range handlers {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/entity", tt.create)
			router.PATCH("/entity/:id", tt.update)
			for _, method := range []string{"POST", "PATCH"} {
				target := "/entity"
				if method == "PATCH" {
					target += "/1"
				}
				request := httptest.NewRequest(method, target, strings.NewReader(" \n null \t"))
				request.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				// A service call would panic: none is configured for this test.
				router.ServeHTTP(recorder, request)
				if recorder.Code != 400 {
					t.Fatalf("%s accepted null: %d %s", method, recorder.Code, recorder.Body)
				}
			}
		})
	}
}

func TestBindJSONObjectRejectsInvalidShapesAndOversizedBodies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	bodies := []string{
		"", "null", "[]", `"string"`, "123", "true",
		"{invalid", "{}{}", `{"nome":"` + strings.Repeat("a", 1<<20) + `"}`,
	}
	for _, body := range bodies {
		router := gin.New()
		router.POST("/", func(ctx *gin.Context) {
			var target struct {
				Nome string `json:"nome"`
			}
			if err := bindJSONObject(ctx, &target); err == nil {
				t.Error("accepted invalid or oversized body")
			}
			ctx.Status(400)
		})
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/", strings.NewReader(body)))
	}
	router := gin.New()
	router.POST("/", func(ctx *gin.Context) {
		var target struct {
			Nome string `json:"nome"`
		}
		if err := bindJSONObject(ctx, &target); err != nil {
			t.Fatal(err)
		}
		if target.Nome != "Nome válido" {
			t.Fatal("valid payload lost its data")
		}
		ctx.Status(200)
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("POST", "/", strings.NewReader(`{"nome":"Nome válido"}`)))
	if recorder.Code != 200 {
		t.Fatal("valid JSON object rejected")
	}
}
