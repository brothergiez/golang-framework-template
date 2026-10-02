package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	appuser "github.com/lion/go-framework/internal/application/user"
	"github.com/lion/go-framework/internal/infrastructure/logger"
	"github.com/lion/go-framework/internal/shared/utils"
)

type DefaultHandler struct {
	users *appuser.Service
	log   *slog.Logger
}

func NewDefaultHandler(users *appuser.Service, log *slog.Logger) *DefaultHandler {
	return &DefaultHandler{users: users, log: log}
}

type createUserEvent struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (h *DefaultHandler) Handle(ctx context.Context, event any) (any, error) {
	rid := utils.NewID()
	ctx = logger.ContextWithRequestID(ctx, rid)
	log := logger.WithContext(ctx, h.log)

	raw, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	var ev createUserEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, fmt.Errorf("invalid event: %w", err)
	}
	if ev.Action == "" {
		ev.Action = "user.create"
	}
	if ev.Action != "user.create" {
		return map[string]any{"ok": true, "skipped": true, "action": ev.Action}, nil
	}

	u, err := h.users.Create(ctx, appuser.CreateUserInput{
		Name:  ev.Name,
		Email: ev.Email,
	})
	if err != nil {
		log.Error("lambda create user failed", "error", err)
		return nil, err
	}
	log.Info("user created from lambda", "user_id", u.ID)
	return map[string]any{
		"ok":      true,
		"user_id": u.ID,
		"email":   u.Email,
	}, nil
}
