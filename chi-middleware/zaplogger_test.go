package chi_middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func TestZapLoggerEntry_Write(t *testing.T) {
	log, _ := zap.NewProduction()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com", nil)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "foobar")
	req = req.WithContext(ctx)
	entry := NewZapLogEntry(log, req)
	entry.Write(200, 100, http.Header{"User-Agent": []string{"gopkg/zaplog"}}, 200*time.Millisecond, nil)
}
