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
		method    string
		path      string
		query     string
		userAgent string
		ip        string
	}
)

func NewZapLogEntry(logger *zap.Logger, r *http.Request) middleware.LogEntry {
	return ZapLoggerEntry{
		logger:    logger,
		requestID: middleware.GetReqID(r.Context()),
		method:    r.Method,
		path:      r.URL.Path,
		query:     r.URL.RawQuery,
		userAgent: r.Header.Get("User-Agent"),
		ip:        r.RemoteAddr,
	}
}

func (zl ZapLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	zl.logger.Info("request complete",
		zap.String("request_method", zl.method),
		zap.String("request_path", zl.path),
		zap.String("request_query", zl.query),
		zap.String("user_agent", zl.userAgent),
		zap.String("client_ip", zl.ip),
		zap.Int("response_status_code", status),
		zap.Int("response_bytes_length", bytes),
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

func RequestZapLog(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := NewZapLogEntry(logger, r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(start), nil)
			}()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}
