package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const reqIDLogKey = "requestID"

type Logger struct {
	zap *zap.Logger
}

func NewLogger() *Logger {
	l, _ := zap.NewProduction(zap.AddCallerSkip(1))

	return &Logger{zap: l}
}

func NewDebugLogger() *Logger {
	l, _ := zap.NewDevelopment(zap.AddCallerSkip(1))

	return &Logger{zap: l}
}

func NewTestLogger() *Logger {
	l := zap.NewNop()

	return &Logger{zap: l}
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	reqID := GetRequestID(ctx)
	fields = append(fields, zap.String(reqIDLogKey, reqID))
	l.zap.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	reqID := GetRequestID(ctx)
	fields = append(fields, zap.String(reqIDLogKey, reqID))
	l.zap.Error(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	reqID := GetRequestID(ctx)
	fields = append(fields, zap.String(reqIDLogKey, reqID))
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) Sync() {
	_ = l.zap.Sync()
}

func (l *Logger) With(fields ...zap.Field) {
	newZap := l.zap.With(fields...)
	l.zap = newZap
}

func Middleware(log *Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := r.Context()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			status := ww.Status()
			duration := time.Since(start)

			fields := []zapcore.Field{
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.String("request", r.RequestURI),
				zap.String("method", r.Method),
			}

			log.Info(ctx, http.StatusText(status), fields...)
		}

		return http.HandlerFunc(fn)
	}
}

func GetRequestID(ctx context.Context) string {
	reqID, ok := ctx.Value(middleware.RequestIDKey).(string)
	if ok {
		return reqID
	}

	return ""
}
