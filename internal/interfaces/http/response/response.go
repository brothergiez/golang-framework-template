package httpresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lion/go-framework/internal/infrastructure/logger"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
	"github.com/lion/go-framework/internal/shared/response"
)

func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, response.Success(data, meta(c)))
}

func Created(c *gin.Context, data any) {
	JSON(c, http.StatusCreated, response.Success(data, meta(c)))
}

func OKWithMeta(c *gin.Context, data any, page, limit int, total int64) {
	m := meta(c)
	m.Page = page
	m.Limit = limit
	m.Total = total
	JSON(c, http.StatusOK, response.Success(data, m))
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, err error) {
	if appErr, ok := apperrors.AsAppError(err); ok {
		JSON(c, appErr.HTTPStatus, response.Fail(appErr.Code, appErr.Message, appErr.Details, meta(c)))
		return
	}
	JSON(c, http.StatusInternalServerError, response.Fail("INTERNAL_ERROR", "internal server error", nil, meta(c)))
}

func JSON(c *gin.Context, status int, body response.Envelope) {
	c.JSON(status, body)
}

func meta(c *gin.Context) *response.Meta {
	return &response.Meta{
		RequestID: logger.RequestIDFromContext(c.Request.Context()),
	}
}
