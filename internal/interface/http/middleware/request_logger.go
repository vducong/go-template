package middleware

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gotemplate/pkg/httprespwrit"
	"gotemplate/pkg/lg"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

func RequestLogger(log lg.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := httprespwrit.New(w)

			var reqBody bytes.Buffer
			r.Body = io.NopCloser(io.TeeReader(r.Body, &reqBody))

			var respBody bytes.Buffer
			ww.Tee(&respBody)

			start := time.Now()

			defer func() {
				var fields []lg.Field

				if rcv := recover(); rcv != nil {
					if ww.Status() == 0 && r.Header.Get("Connection") != "Upgrade" {
						ww.WriteHeader(http.StatusInternalServerError)
					}

					if rcv == http.ErrAbortHandler {
						defer panic(rcv)
					}

					fields = append(fields, lg.Str(lg.ErrorFieldName, fmt.Sprintf("panic: %+v", rcv)))

					pc := make([]uintptr, 10)
					n := runtime.Callers(3, pc)
					pc = pc[:n]
					stacks := panicStacksFromProgramCounter(pc)
					fields = append(fields, lg.Str(lg.StackFieldName, strings.Join(stacks, "\n")))
				}

				duration := time.Since(start)
				statusCode := ww.Status()
				if statusCode == 0 {
					statusCode = http.StatusOK
				}

				ctx := r.Context()

				fields = append(fields,
					lg.Str("request_id", chimw.GetReqID(ctx)),
					lg.Str("method", r.Method),
					lg.Str("path", r.URL.Path),
					lg.Str("ip", clientIPFromRequest(r)),
					lg.Str("http.request.body", logBody(&reqBody)),
					lg.Int("http.request.body.size", int(r.ContentLength)),
					lg.Str("user_agent", r.UserAgent()),
					lg.Str("referer", r.Referer()),
					lg.Int("status", statusCode),
					lg.Str("latency", duration.String()),
					lg.Str("http.response.body", logBody(&respBody)),
					lg.Int("http.response.body.size", ww.BytesWritten()),
				)

				if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
					spanCtx := span.SpanContext()
					fields = append(fields,
						lg.Str("trace_id", spanCtx.TraceID().String()),
						lg.Str("span_id", spanCtx.SpanID().String()),
					)
				}

				if err := ctx.Err(); errors.Is(err, context.Canceled) {
					fields = append(fields, lg.Str(lg.ErrorFieldName, fmt.Sprintf("client aborted: %v", err)))
				}

				logLevel := logLevelFromStatusCode(statusCode)
				msg := fmt.Sprintf("%s %s => HTTP %v (%v)", r.Method, r.URL, statusCode, duration)
				log.CtxLog(ctx, logLevel, msg, fields...)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func panicStacksFromProgramCounter(pc []uintptr) []string {
	frames := runtime.CallersFrames(pc)
	var stacks []string
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if !strings.Contains(frame.File, "runtime/panic.go") {
			stacks = append(stacks, fmt.Sprintf("%s:%d", frame.File, frame.Line))
		}
	}
	return stacks
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

func logBody(body *bytes.Buffer) string {
	if body.Len() == 0 {
		return ""
	}
	return body.String()
}
