package config_test

import (
	"os"
	"testing"

	"github.com/lion/go-framework/internal/config"
)

func TestLoadDefaultMemoryAPI(t *testing.T) {
	t.Setenv("SERVICE_TYPE", "api")
	t.Setenv("DATABASE_DRIVER", "memory")
	t.Setenv("APP_NAME", "test-service")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Service.Type != config.ServiceTypeAPI {
		t.Fatalf("unexpected type: %s", cfg.Service.Type)
	}
	if cfg.Database.Driver != "memory" {
		t.Fatalf("unexpected driver: %s", cfg.Database.Driver)
	}
}

func TestInvalidServiceType(t *testing.T) {
	t.Setenv("SERVICE_TYPE", "foobar")
	t.Setenv("DATABASE_DRIVER", "memory")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error")
	}
	_ = os.Unsetenv("SERVICE_TYPE")
}
