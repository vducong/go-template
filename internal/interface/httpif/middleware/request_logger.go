package middleware

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gotemplate/pkg/lg"
	"io"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

// requestLogBodyLimit caps each body logged for a failed request.
const requestLogBodyLimit = 4 << 10

func RequestLogger(log lg.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			reqBody := &cappedBuffer{limit: requestLogBodyLimit}
			r.Body = io.NopCloser(io.TeeReader(r.Body, reqBody))

			respBody := &cappedBuffer{limit: requestLogBodyLimit}
			ww.Tee(respBody)

			start := time.Now()
			defer logCompletedRequest(log, r, ww, reqBody, respBody, start)

			next.ServeHTTP(ww, r)
		})
	}
}

func logCompletedRequest(
	log lg.Logger,
	r *http.Request,
	ww chimw.WrapResponseWriter,
	reqBody, respBody *cappedBuffer,
	start time.Time,
) {
	duration := time.Since(start)
	statusCode := ww.Status()
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	ctx := r.Context()
	fields := requestFields(r, ww, reqBody, respBody, statusCode, duration)
	fields = append(fields, traceFields(ctx)...)
	fields = append(fields, clientAbortFields(ctx)...)

	msg := fmt.Sprintf("%s %s => HTTP %v (%v)", r.Method, r.URL, statusCode, duration)
	log.CtxLog(ctx, logLevelFromStatusCode(statusCode), msg, fields...)
}

func requestFields(
	r *http.Request,
	ww chimw.WrapResponseWriter,
	reqBody, respBody *cappedBuffer,
	statusCode int,
	duration time.Duration,
) []lg.Field {
	fields := []lg.Field{
		lg.Str("request_id", chimw.GetReqID(r.Context())),
		lg.Str("method", r.Method),
		lg.Str("path", r.URL.Path),
		lg.Str("ip", clientIPFromRequest(r)),
		lg.Int("http.request.body.size", int(r.ContentLength)),
		lg.Str("user_agent", r.UserAgent()),
		lg.Str("referer", r.Referer()),
		lg.Int("status", statusCode),
		lg.Str("latency", duration.String()),
		lg.Int("http.response.body.size", ww.BytesWritten()),
	}
	// Bodies only for failed requests: a successful one can carry personal or money data.
	if statusCode >= http.StatusBadRequest {
		fields = append(fields,
			lg.Str("http.request.body", reqBody.String()),
			lg.Str("http.response.body", respBody.String()),
		)
	}
	return fields
}

func traceFields(ctx context.Context) []lg.Field {
	span := trace.SpanFromContext(ctx)
	// An unsampled trace is never exported, so its trace_id would point at nothing.
	if !span.SpanContext().IsSampled() {
		return nil
	}
	spanCtx := span.SpanContext()
	return []lg.Field{
		lg.Str("trace_id", spanCtx.TraceID().String()),
		lg.Str("span_id", spanCtx.SpanID().String()),
	}
}

func clientAbortFields(ctx context.Context) []lg.Field {
	err := ctx.Err()
	if !errors.Is(err, context.Canceled) {
		return nil
	}
	return []lg.Field{
		lg.Str(lg.ErrorFieldName, fmt.Sprintf("client aborted: %v", err)),
	}
}

func logLevelFromStatusCode(statusCode int) lg.LogLevel {
	if statusCode >= http.StatusInternalServerError {
		return lg.LogLevelError
	}
	if statusCode >= http.StatusBadRequest {
		return lg.LogLevelWarn
	}
	return lg.LogLevelInfo
}

func clientIPFromRequest(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

// cappedBuffer keeps the first limit bytes and reports every write as complete,
// so the TeeReader in front of it never cuts the handler's read short.
type cappedBuffer struct {
	buf   bytes.Buffer
	limit int
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room > 0 {
		_, _ = b.buf.Write(p[:min(room, len(p))])
	}
	return len(p), nil
}

func (b *cappedBuffer) String() string {
	return b.buf.String()
}
