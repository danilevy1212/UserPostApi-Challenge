package logger

// TODO  TEST!!!

import (
	"bytes"
	"context"
	"io"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ctxKeyLogger struct{}

var loggerKey = ctxKeyLogger{}

func FromContext(ctx context.Context) *zerolog.Logger {
	l, ok := ctx.Value(loggerKey).(*zerolog.Logger)
	if !ok {
		panic("logger not found in context, did you forget to use logger middleware?")
	}

	return l
}

// Captures the gin.Response and logs it before sending it
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewMiddleware(baseLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		reqLogger := baseLogger.With().
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.Request.UserAgent()).
			Logger()

		// Preparing to capture the response buffer
		respBuf := new(bytes.Buffer)
		writer := &bodyWriter{
			ResponseWriter: ctx.Writer,
			body:           respBuf,
		}
		ctx.Writer = writer

		// Inject logger
		ctxWithLogger := withContext(ctx.Request.Context(), &reqLogger)
		ctx.Request = ctx.Request.WithContext(ctxWithLogger)

		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		reqLogger.Info().
			Str("request_body", string(requestBody)).
			Interface("request_headers", ctx.Request.Header).
			Msg("Processing request")

		// Proceed with request
		ctx.Next()

		// Final log
		reqLogger.Info().
			Int("status", ctx.Writer.Status()).
			Str("response_body", respBuf.String()).
			Interface("response_headers", ctx.Writer.Header()).
			Dur("duration", time.Since(start)).
			Msg("response sent")
	}
}

func withContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
