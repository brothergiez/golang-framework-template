package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	appuser "github.com/lion/go-framework/internal/application/user"
	"github.com/lion/go-framework/internal/app/cron"
	applambda "github.com/lion/go-framework/internal/app/lambda"
	"github.com/lion/go-framework/internal/app/server"
	"github.com/lion/go-framework/internal/app/worker"
	"github.com/lion/go-framework/internal/config"
	"github.com/lion/go-framework/internal/infrastructure/cache"
	rediscache "github.com/lion/go-framework/internal/infrastructure/cache/redis"
	"github.com/lion/go-framework/internal/infrastructure/database"
	"github.com/lion/go-framework/internal/infrastructure/logger"
	"github.com/lion/go-framework/internal/interfaces/cron/jobs"
	"github.com/lion/go-framework/internal/interfaces/http/handler"
	"github.com/lion/go-framework/internal/interfaces/http/router"
	lambdahandlers "github.com/lion/go-framework/internal/interfaces/lambda/handlers"
	"github.com/lion/go-framework/internal/interfaces/worker/consumers"
)

type Application struct {
	Config     *config.Config
	Logger     *slog.Logger
	Cache      cache.Cache
	DB         *database.Resources
	UserService *appuser.Service

	APIServer   *server.Server
	Scheduler   *cron.IntervalScheduler
	Worker      *worker.Runner
	WorkerQueue *worker.InMemoryQueue
	Lambda      *applambda.Adapter
}

func New() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.Logger).With(
		"service", cfg.App.Name,
		"env", cfg.App.Env,
		"service_type", cfg.Service.Type,
	)

	dbRes, err := database.NewResources(cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	var cacheStore cache.Cache = cache.NewNoop()
	if cfg.Redis.Enabled {
		rc, err := rediscache.New(cfg.Redis)
		if err != nil {
			_ = dbRes.Close()
			return nil, fmt.Errorf("init redis: %w", err)
		}
		cacheStore = rc
	}

	userService := appuser.NewService(dbRes.UserRepo)

	app := &Application{
		Config:      cfg,
		Logger:      log,
		Cache:       cacheStore,
		DB:          dbRes,
		UserService: userService,
	}

	switch cfg.Service.Type {
	case config.ServiceTypeAPI:
		userHandler := handler.NewUserHandler(userService)
		healthHandler := handler.NewHealthHandler(dbRes.Ping, cacheStore, cfg.Redis.Enabled)
		engine := router.New(router.Dependencies{
			Config: cfg,
			Logger: log,
			Users:  userHandler,
			Health: healthHandler,
		})
		app.APIServer = server.New(cfg, engine, log)
	case config.ServiceTypeCron:
		scheduler := cron.NewIntervalScheduler()
		job := jobs.NewCleanupUsersJob(userService, log)
		if err := scheduler.Register("1m", job); err != nil {
			_ = app.Close()
			return nil, err
		}
		app.Scheduler = scheduler
	case config.ServiceTypeWorker:
		queue := worker.NewInMemoryQueue(100)
		consumer := consumers.NewCreateUserConsumer(userService, log)
		app.WorkerQueue = queue
		app.Worker = worker.NewRunner(consumer, queue, cfg.Worker.Concurrency, cfg.Worker.PollInterval)
	case config.ServiceTypeLambda:
		h := lambdahandlers.NewDefaultHandler(userService, log)
		app.Lambda = applambda.NewAdapter(h)
	default:
		_ = app.Close()
		return nil, fmt.Errorf("unsupported service type: %s", cfg.Service.Type)
	}

	return app, nil
}

func (a *Application) Run(ctx context.Context) error {
	switch a.Config.Service.Type {
	case config.ServiceTypeAPI:
		errCh := make(chan error, 1)
		go func() {
			errCh <- a.APIServer.Start()
		}()
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), a.Config.App.ShutdownTimeout)
			defer cancel()
			_ = a.APIServer.Stop(shutdownCtx)
			return a.Close()
		case err := <-errCh:
			_ = a.Close()
			return err
		}
	case config.ServiceTypeCron:
		errCh := make(chan error, 1)
		go func() {
			errCh <- a.Scheduler.Start(ctx)
		}()
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), a.Config.App.ShutdownTimeout)
			defer cancel()
			_ = a.Scheduler.Stop(shutdownCtx)
			return a.Close()
		case err := <-errCh:
			_ = a.Close()
			return err
		}
	case config.ServiceTypeWorker:
		// Demo message so worker does something when started locally.
		go func() {
			time.Sleep(time.Second)
			a.WorkerQueue.Publish(worker.Message{
				ID:      "demo-1",
				Type:    "user.create",
				Payload: []byte(`{"name":"Worker User","email":"worker@example.com"}`),
			})
		}()
		errCh := make(chan error, 1)
		go func() {
			errCh <- a.Worker.Start(ctx)
		}()
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), a.Config.App.ShutdownTimeout)
			defer cancel()
			_ = a.Worker.Stop(shutdownCtx)
			return a.Close()
		case err := <-errCh:
			_ = a.Close()
			return err
		}
	case config.ServiceTypeLambda:
		// Lambda entrypoints invoke adapters directly; Run keeps process ready for local demo.
		a.Logger.Info("lambda adapter ready; waiting for invoke or signal")
		<-ctx.Done()
		return a.Close()
	default:
		return fmt.Errorf("unsupported service type: %s", a.Config.Service.Type)
	}
}

func (a *Application) Close() error {
	var firstErr error
	if a.Cache != nil {
		if err := a.Cache.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if a.DB != nil && a.DB.Close != nil {
		if err := a.DB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
