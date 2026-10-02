package user_test

import (
	"context"
	"testing"

	appuser "github.com/lion/go-framework/internal/application/user"
	"github.com/lion/go-framework/internal/infrastructure/database/memory"
	"github.com/lion/go-framework/internal/shared/pagination"
)

func TestUserServiceCRUD(t *testing.T) {
	repo := memory.NewUserRepository()
	svc := appuser.NewService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, appuser.CreateUserInput{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected id")
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Email != "alice@example.com" {
		t.Fatalf("unexpected email: %s", got.Email)
	}

	_, err = svc.Create(ctx, appuser.CreateUserInput{
		Name:  "Alice 2",
		Email: "alice@example.com",
	})
	if err == nil {
		t.Fatal("expected conflict")
	}

	list, err := svc.List(ctx, pagination.Params{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Total != 1 {
		t.Fatalf("expected total 1, got %d", list.Total)
	}

	updated, err := svc.Update(ctx, created.ID, appuser.UpdateUserInput{Name: "Alice Updated"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Fatalf("unexpected name: %s", updated.Name)
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.GetByID(ctx, created.ID); err == nil {
		t.Fatal("expected not found after delete")
	}
}
