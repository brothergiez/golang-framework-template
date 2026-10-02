package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	appuser "github.com/lion/go-framework/internal/application/user"
	domainuser "github.com/lion/go-framework/internal/domain/user"
	"github.com/lion/go-framework/internal/interfaces/http/request"
	httpresponse "github.com/lion/go-framework/internal/interfaces/http/response"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
	"github.com/lion/go-framework/internal/shared/pagination"
	"github.com/lion/go-framework/internal/shared/validator"
)

type UserHandler struct {
	service *appuser.Service
}

func NewUserHandler(service *appuser.Service) *UserHandler {
	return &UserHandler{service: service}
}

type userDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperrors.BadRequest("invalid JSON body"))
		return
	}
	if err := validator.Struct(req); err != nil {
		httpresponse.Error(c, err)
		return
	}
	u, err := h.service.Create(c.Request.Context(), appuser.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	httpresponse.Created(c, toDTO(u))
}

func (h *UserHandler) GetByID(c *gin.Context) {
	u, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	httpresponse.OK(c, toDTO(u))
}

func (h *UserHandler) List(c *gin.Context) {
	params := pagination.Parse(c.Query("page"), c.Query("limit"), 20, 100)
	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	items := make([]userDTO, 0, len(result.Items))
	for _, u := range result.Items {
		items = append(items, toDTO(u))
	}
	httpresponse.OKWithMeta(c, items, result.Page, result.Limit, result.Total)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperrors.BadRequest("invalid JSON body"))
		return
	}
	if err := validator.Struct(req); err != nil {
		httpresponse.Error(c, err)
		return
	}
	u, err := h.service.Update(c.Request.Context(), c.Param("id"), appuser.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	httpresponse.OK(c, toDTO(u))
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		httpresponse.Error(c, err)
		return
	}
	httpresponse.NoContent(c)
}

func toDTO(u *domainuser.User) userDTO {
	return userDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
