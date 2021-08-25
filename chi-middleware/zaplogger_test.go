package chi_middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func TestZapLoggerEntry_Write(t *testing.T) {
	logger, _ := zap.NewProduction()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(RequestZapLog(logger))
	r.Get("/foobar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("foobar"))
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:19988",
		Handler: r,
	}
	go srv.ListenAndServe()

	time.Sleep(1 * time.Second)
	http.Get("http://127.0.0.1:19988/foobar")
	srv.Shutdown(context.TODO())
}
