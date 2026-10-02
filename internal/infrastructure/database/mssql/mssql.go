package mssql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/lion/go-framework/internal/config"
	domainuser "github.com/lion/go-framework/internal/domain/user"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
)

func Open(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", cfg.DSN())
	if err != nil {
		return nil, err
	}
	configurePool(db, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='users' AND xtype='U')
CREATE TABLE users (
    id NVARCHAR(64) PRIMARY KEY,
    name NVARCHAR(255) NOT NULL,
    email NVARCHAR(255) NOT NULL UNIQUE,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2 NOT NULL
)`)
	return err
}

func (r *UserRepository) Create(ctx context.Context, user *domainuser.User) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO users (id, name, email, created_at, updated_at)
VALUES (@p1, @p2, @p3, @p4, @p5)`,
		user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	return r.scanOne(ctx, `SELECT id, name, email, created_at, updated_at FROM users WHERE id = @p1`, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	return r.scanOne(ctx, `SELECT id, name, email, created_at, updated_at FROM users WHERE email = @p1`, email)
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*domainuser.User, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT id, name, email, created_at, updated_at
FROM users
ORDER BY created_at DESC
OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY`, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*domainuser.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, u)
	}
	return items, total, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *domainuser.User) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE users SET name = @p2, email = @p3, updated_at = @p4 WHERE id = @p1`,
		user.ID, user.Name, user.Email, user.UpdatedAt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = @p1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *UserRepository) scanOne(ctx context.Context, query string, args ...any) (*domainuser.User, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	return u, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*domainuser.User, error) {
	var u domainuser.User
	if err := s.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func configurePool(db *sql.DB, cfg config.DatabaseConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}
