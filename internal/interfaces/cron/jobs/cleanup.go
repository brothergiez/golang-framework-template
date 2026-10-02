package jobs

import (
	"context"
	"log/slog"

	appuser "github.com/lion/go-framework/internal/application/user"
	"github.com/lion/go-framework/internal/infrastructure/logger"
	"github.com/lion/go-framework/internal/shared/pagination"
	"github.com/lion/go-framework/internal/shared/utils"
)

type CleanupUsersJob struct {
	users *appuser.Service
	log   *slog.Logger
}

func NewCleanupUsersJob(users *appuser.Service, log *slog.Logger) *CleanupUsersJob {
	return &CleanupUsersJob{users: users, log: log}
}

func (j *CleanupUsersJob) Name() string { return "cleanup_users" }

func (j *CleanupUsersJob) Run(ctx context.Context) error {
	rid := utils.NewID()
	ctx = logger.ContextWithRequestID(ctx, rid)
	log := logger.WithContext(ctx, j.log)

	result, err := j.users.List(ctx, pagination.Params{Page: 1, Limit: 1})
	if err != nil {
		log.Error("cleanup job failed", "error", err)
		return err
	}
	log.Info("cleanup job completed", "user_count_sample_total", result.Total)
	return nil
}
