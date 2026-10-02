package consumers

import (
	"context"
	"encoding/json"
	"log/slog"

	appuser "github.com/lion/go-framework/internal/application/user"
	"github.com/lion/go-framework/internal/app/worker"
	"github.com/lion/go-framework/internal/infrastructure/logger"
	"github.com/lion/go-framework/internal/shared/utils"
)

type CreateUserConsumer struct {
	users *appuser.Service
	log   *slog.Logger
}

func NewCreateUserConsumer(users *appuser.Service, log *slog.Logger) *CreateUserConsumer {
	return &CreateUserConsumer{users: users, log: log}
}

type createUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *CreateUserConsumer) Handle(ctx context.Context, message worker.Message) error {
	rid := message.ID
	if rid == "" {
		rid = utils.NewID()
	}
	ctx = logger.ContextWithRequestID(ctx, rid)
	log := logger.WithContext(ctx, c.log)

	if message.Type != "user.create" {
		log.Warn("unsupported message type", "type", message.Type)
		return nil
	}

	var payload createUserPayload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		log.Error("invalid payload", "error", err)
		return err
	}

	u, err := c.users.Create(ctx, appuser.CreateUserInput{
		Name:  payload.Name,
		Email: payload.Email,
	})
	if err != nil {
		log.Error("create user failed", "error", err)
		return err
	}
	log.Info("user created from worker", "user_id", u.ID)
	return nil
}
