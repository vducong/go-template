package middleware

import (
	"context"
	"gotemplate/pkg/lg"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestRequestLogger_OmitsBodiesOnSuccess(t *testing.T) {
	logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	serve(t, logger, http.StatusOK, `{"password":"hunter2"}`, `{"balance":12345}`)

	fields := logger.only(t)
	for _, key := range []string{"http.request.body", "http.response.body"} {
		if _, ok := fields[key]; ok {
			t.Errorf("%s must not be logged for a 200: %v", key, fields[key])
		}
	}
	if fields["http.request.body.size"] != 22 || fields["http.response.body.size"] != 17 {
		t.Errorf("body sizes must still be logged: %v", fields)
	}
}

func TestRequestLogger_LogsBodiesOnFailure(t *testing.T) {
	logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	serve(t, logger, http.StatusBadRequest, `{"amount":-1}`, `{"error":"bad amount"}`)

	fields := logger.only(t)
	if fields["http.request.body"] != `{"amount":-1}` || fields["http.response.body"] != `{"error":"bad amount"}` {
		t.Errorf("bodies must be logged for a 400: %v", fields)
	}
}

func TestRequestLogger_CapsLoggedBodiesButNotTheRequest(t *testing.T) {
	logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	big := strings.Repeat("x", 10<<10)
	got, rec := serve(t, logger, http.StatusInternalServerError, big, big)

	if len(got) != len(big) || rec.Body.Len() != len(big) {
		t.Fatalf("handler read %d and client got %d bytes, want %d each", len(got), rec.Body.Len(), len(big))
	}
	fields := logger.only(t)
	for _, key := range []string{"http.request.body", "http.response.body"} {
		if body, _ := fields[key].(string); len(body) != requestLogBodyLimit {
			t.Errorf("%s logged %d bytes, want %d", key, len(body), requestLogBodyLimit)
		}
	}
}

// serve runs one request through RequestLogger and returns what the handler read from the body.
func serve(t *testing.T, logger lg.Logger, status int, reqBody, respBody string) (string, *httptest.ResponseRecorder) {
	t.Helper()
	var read []byte
	h := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if read, err = io.ReadAll(r.Body); err != nil {
			t.Error(err)
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, respBody)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(reqBody)))
	return string(read), rec
}

// fieldsLogger records the fields of each request log line RequestLogger writes through CtxLog.
type fieldsLogger struct {
	lg.Logger
	mu    sync.Mutex
	lines []map[string]any
}

func (l *fieldsLogger) CtxLog(_ context.Context, _ lg.LogLevel, _ string, fields ...lg.Field) {
	line := make(map[string]any, len(fields))
	for _, f := range fields {
		line[f.Key] = f.Value
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, line)
}

func (l *fieldsLogger) only(t *testing.T) map[string]any {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.lines) != 1 {
		t.Fatalf("got %d request log lines, want 1", len(l.lines))
	}
	return l.lines[0]
}
