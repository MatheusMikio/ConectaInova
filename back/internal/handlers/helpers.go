package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"conectainova/internal/domain/models"
	query "conectainova/internal/domain/models/query"
	"conectainova/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func bindJSONObject(ctx *gin.Context, target any) error {
	reader := http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1<<20)
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 || body[0] != '{' {
		return errors.New("o corpo da requisição deve ser um objeto JSON")
	}
	if err := json.Unmarshal(body, target); err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(target)
}

func getID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("id", "ID inválido."))
		return 0, false
	}
	return uint(id), true
}

func getPagination(ctx *gin.Context) (query.PageQuery, bool) {
	var pagination query.PageQuery
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("pagination", "Página ou tamanho inválido."))
		return pagination, false
	}
	return query.NewPaginationQuery(pagination.Page, pagination.Size), true
}

func sendServiceError(ctx *gin.Context, err error) {
	var validation *models.ErrorMessage
	switch {
	case errors.As(err, &validation):
		response.SendError(ctx, http.StatusBadRequest, validation)
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.SendError(ctx, http.StatusNotFound, models.NewErrorMessage("id", "Recurso não encontrado."))
	case errors.Is(err, gorm.ErrDuplicatedKey), errors.Is(err, gorm.ErrForeignKeyViolated):
		response.SendError(ctx, http.StatusConflict, models.NewErrorMessage("data", "Conflito com os dados cadastrados."))
	default:
		_ = ctx.Error(err)
		response.SendError(ctx, http.StatusInternalServerError, models.NewErrorMessage("System", "Erro inesperado."))
	}
}
