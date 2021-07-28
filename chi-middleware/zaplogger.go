package chi_middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type (
	ZapLoggerEntry struct {
		logger    *zap.Logger
		requestID string
	}
)

var (
	_ middleware.LogEntry = (*ZapLoggerEntry)(nil)
)

func NewZapLogEntry(logger *zap.Logger, r *http.Request) middleware.LogEntry {
	return ZapLoggerEntry{
		logger:    logger,
		requestID: middleware.GetReqID(r.Context()),
	}
}

func (zl ZapLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	zl.logger.Info("request complete",
		zap.Int("response_status_code", status),
		zap.Int("response_bytes_length", bytes),
		zap.String("user_agent", header.Get("User-Agent")),
		zap.String("elapsed", elapsed.String()),
		zap.String("request_id", zl.requestID),
	)
}

func (zl ZapLoggerEntry) Panic(v interface{}, stack []byte) {
	zl.logger.Error("broken request",
		zap.Any("panic", v),
		zap.ByteString("stack", stack),
		zap.String("request_id", zl.requestID),
	)
}

func RequestZapLog(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := NewZapLogEntry(logger, r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))

			entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(start), nil)
		}
		return http.HandlerFunc(fn)
	}
}
