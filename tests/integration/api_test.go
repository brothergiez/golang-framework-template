package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lion/go-framework/internal/app/bootstrap"
	"github.com/lion/go-framework/internal/interfaces/http/handler"
	"github.com/lion/go-framework/internal/interfaces/http/router"
)

func TestAPICreateAndGetUser(t *testing.T) {
	_ = os.Setenv("SERVICE_TYPE", "api")
	_ = os.Setenv("DATABASE_DRIVER", "memory")
	_ = os.Setenv("APP_NAME", "integration-test")
	_ = os.Setenv("REDIS_ENABLED", "false")

	app, err := bootstrap.New()
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	defer func() { _ = app.Close() }()

	engine := router.New(router.Dependencies{
		Config: app.Config,
		Logger: app.Logger,
		Users:  handler.NewUserHandler(app.UserService),
		Health: handler.NewHealthHandler(app.DB.Ping, app.Cache, false),
	})

	body := map[string]string{"name": "Bob", "email": "bob@example.com"}
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", w.Code, w.Body.String())
	}

	var created struct {
		Success bool `json:"success"`
		Data    struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.Data.ID == "" {
		t.Fatal("missing id")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/"+created.Data.ID, nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get status=%d body=%s", w.Code, w.Body.String())
	}
}
