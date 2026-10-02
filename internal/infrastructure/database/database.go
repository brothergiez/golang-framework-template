package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lion/go-framework/internal/config"
	domainuser "github.com/lion/go-framework/internal/domain/user"
	"github.com/lion/go-framework/internal/infrastructure/database/memory"
	"github.com/lion/go-framework/internal/infrastructure/database/mongo"
	"github.com/lion/go-framework/internal/infrastructure/database/mssql"
	"github.com/lion/go-framework/internal/infrastructure/database/mysql"
	"github.com/lion/go-framework/internal/infrastructure/database/postgres"
)

type Resources struct {
	SQL      *sql.DB
	Mongo    *mongo.Client
	UserRepo domainuser.Repository
	Ping     func(ctx context.Context) error
	Close    func() error
}

func NewResources(cfg *config.Config) (*Resources, error) {
	switch cfg.Database.Driver {
	case "memory":
		repo := memory.NewUserRepository()
		return &Resources{
			UserRepo: repo,
			Ping:     func(context.Context) error { return nil },
			Close:    func() error { return nil },
		}, nil
	case "postgres":
		db, err := postgres.Open(cfg.Database)
		if err != nil {
			return nil, err
		}
		repo := postgres.NewUserRepository(db)
		if err := repo.Migrate(context.Background()); err != nil {
			_ = db.Close()
			return nil, err
		}
		return &Resources{
			SQL:      db,
			UserRepo: repo,
			Ping:     db.PingContext,
			Close:    db.Close,
		}, nil
	case "mysql", "mariadb":
		db, err := mysql.Open(cfg.Database)
		if err != nil {
			return nil, err
		}
		repo := mysql.NewUserRepository(db)
		if err := repo.Migrate(context.Background()); err != nil {
			_ = db.Close()
			return nil, err
		}
		return &Resources{
			SQL:      db,
			UserRepo: repo,
			Ping:     db.PingContext,
			Close:    db.Close,
		}, nil
	case "mssql":
		db, err := mssql.Open(cfg.Database)
		if err != nil {
			return nil, err
		}
		repo := mssql.NewUserRepository(db)
		if err := repo.Migrate(context.Background()); err != nil {
			_ = db.Close()
			return nil, err
		}
		return &Resources{
			SQL:      db,
			UserRepo: repo,
			Ping:     db.PingContext,
			Close:    db.Close,
		}, nil
	case "mongo":
		client, err := mongo.Open(cfg.Mongo)
		if err != nil {
			return nil, err
		}
		repo := mongo.NewUserRepository(client, cfg.Mongo.Database)
		return &Resources{
			Mongo:    client,
			UserRepo: repo,
			Ping:     client.Ping,
			Close:    client.Close,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
}

func PingWithTimeout(ping func(context.Context) error, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ping(ctx)
}
