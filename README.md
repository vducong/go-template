# go-template

A bootstrap template for Go microservices featuring dual-stack HTTP (Chi) and gRPC, JWT & API key authentication, structured error boundaries, and OpenTelemetry instrumentation.

---

## 💡 Architecture & Design Goals

This codebase is structured around **Clean Architecture** principles to separate concerns and ensure transport-independent domain logic:

- **Inward Dependencies**: Transport (`internal/interface`) and storage drivers (`internal/repository`) depend inward on business logic (`internal/service`).
- **Error Boundary**: Internal failures are intercepted and translated into secure, safe-for-client error codes without leaking implementation details.
- **Built-in Telemetry**: Metrics (Prometheus), structured logging (zerolog), tracing (OpenTelemetry/Zipkin), and panic recovery are wired at the outer server boundary.
- **Testability**: Interfaces are mockable out of the box using Uber Mock (`go.uber.org/mock`).

---

## ⚡ Quick Start

### 1. Initialize Configuration

```bash
cp configs/example.config.yaml configs/config.yaml
```

_(By default, `configs/config.yaml` is gitignored so local credentials remain safe.)_

### 2. Start the Servers

```bash
make run
```

Both HTTP (default `:8080`) and gRPC servers will start.

### 3. Verify Health Check

```bash
curl -i http://localhost:8080/health
```

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"success":true,"data":"OK"}
```

---

## 🏗️ Codebase Layout

```
.
├── cmd/
│   ├── app/                # Main server entrypoint (cfg.Load("") relies on default fallback)
│   └── cli/                # CLI runner supporting the -config-file flag
├── configs/                # Default configuration files
├── internal/
│   ├── app/                # Initialization, dependency injection, and lifecycles
│   ├── apperr/             # Error domain mapping (service.domain.error)
│   ├── cfg/                # Config structs and loading logic
│   ├── infra/              # Shared infrastructure clients (DB, MQ, event bus, logger, telemetry)
│   ├── interface/          # HTTP & gRPC transport layers (delivery boundary)
│   │   ├── grpcif/         # gRPC server bindings and interceptors
│   │   └── httpif/         # chi Router, middleware, and handlers
│   ├── repository/         # Data storage access (DB, cache, external APIs)
│   └── service/            # Core business rules (pure, transport-agnostic logic)
└── pkg/                    # Reusable packages (auth, config-loaders, server-init, etc.)
```

---

## 🧭 Developer Guides

### 1. Adding an Endpoint

Follow the dependency flow from the repository up to routing:

#### Step A: Define the Storage Interface (`internal/repository/user.go`)

```go
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=mock -typed
package repository

import "context"

type UserStore interface {
    Get(ctx context.Context, id string) (string, error)
}
```

#### Step B: Implement Business Logic (`internal/service/user.go`)

Keep logic transport-agnostic (no HTTP/gRPC specific imports here):

```go
package service

import (
    "context"
    "gotemplate/internal/repository"
)

type UserService struct {
    store repository.UserStore
}

func (s *UserService) GetUserProfile(ctx context.Context, id string) (string, error) {
    return s.store.Get(ctx, id)
}
```

#### Step C: Implement Transport Handler (`internal/interface/httpif/handler/user.go`)

```go
package handler

import (
    "gotemplate/pkg/httprespwrit"
    "net/http"
    "github.com/go-chi/chi/v5"
)

type UserHandler struct {
    ResponseWriter httprespwrit.Writer
    UserService    *service.UserService
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    profile, err := h.UserService.GetUserProfile(r.Context(), userID)
    if err != nil {
        h.ResponseWriter.WriteError(w, r, &httprespwrit.ErrorResponse{Err: err})
        return
    }

    h.ResponseWriter.WriteJSON(w, r, &httprespwrit.JSONResponse{
        StatusCode: http.StatusOK,
        Data:       profile,
    })
}
```

#### Step D: Mount the Route (`internal/interface/httpif/server.go`)

```go
r.Route("/api/v1", func(api chi.Router) {
    api.Use(authenticator.RequireJwt()) // JWT Middleware
    api.Get("/users/{id}", handlers.User.Get)
})
```

---

### 2. Throwing and Propagating Structured Errors

The template translates internal errors into secure, structured API messages matching the `{Service}.{Domain}.{Number}` format.

#### Step A: Define the Code & Map Status (`internal/apperr/codes.go`)

```go
var CodeUserNotFound = Code{Domain: 1, Number: 1}

var errorCodeToStatusCode = map[Code]int{
    CodeUserNotFound: http.StatusNotFound,
}

var errorCodeToKey = map[Code]string{
    CodeUserNotFound: "user.not_found",
}
```

#### Step B: Register the Message Catalog (`internal/apperr/messages.json`)

_(Keys are validated at startup to prevent unmapped codes)_

```json
{
  "gotemplate.user.not_found": "The requested user profile was not found."
}
```

#### Step C: Throw and Response Output

```go
// Return from service or repository:
return "", &apperr.Error{ErrCode: apperr.CodeUserNotFound}
```

The transport-level response writer outputs:

```json
{
  "success": false,
  "error": {
    "code": "12.01.01",
    "key": "gotemplate.user.not_found",
    "detail": "The requested user profile was not found."
  }
}
```

---

### 3. Running and Generating Mocks

Code generation uses **Uber Mock (`go.uber.org/mock`)**.

#### Step A: Regenerate Mocks

Run generation whenever interfaces change:

```bash
make gen
```

_(Runs `go generate ./...` which populates matching `/mock` folders)_

#### Step B: Use Mocks in Tests

```go
func TestUserService_GetUserProfile(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockStore := mock.NewMockUserStore(ctrl)
    mockStore.EXPECT().
        Get(gomock.Any(), "user-123").
        Return("Test User Profile", nil).
        Times(1)

    svc := &UserService{store: mockStore}
    profile, err := svc.GetUserProfile(context.Background(), "user-123")

    assert.NoError(t, err)
    assert.Equal(t, "Test User Profile", profile)
}
```

#### Step C: Execute Tests

```bash
make test
```

---

## ⚙️ Configuration Loading Sequence

Config options are loaded using `pkg/cfld` in a progressive fallback pipeline:

```
┌────────────────────────────────────────────────────────┐
│  1. Explicit Flag (-config-file <path>)                │ (e.g., /configs/config.yaml)
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│  2. CONFIG_FILE_PATH Environment Variable              │
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│  3. Default File (configs/config.yaml)                 │
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│  4. Environment Variable Overrides & Struct Defaults   │ (e.g., HTTP_PORT, LOG_LEVEL)
└────────────────────────────────────────────────────────┘
```

---

## 📊 Telemetry Options

- **Logs (`pkg/lg`)**: Supports `console` (pretty text) and `json` formats via zerolog.
- **Metrics (`pkg/mtr`)**: Standard Prometheus metrics exposed under `GET /metrics`.
- **Traces (`pkg/trc`)**: Standard OpenTelemetry spans exportable to Zipkin, or OTLP (gRPC/HTTP).

---

## 🛠️ Make targets

| Command             | Action                                                                                                        |
| ------------------- | ------------------------------------------------------------------------------------------------------------- |
| `make install`      | Installs system dependencies (`buf`, `golangci-lint`), run `go mod tidy`, and registers git pre-commit hooks. |
| `make build`        | Compiles the binary to `bin/app`.                                                                             |
| `make run`          | Runs `cmd/app` using local fallbacks.                                                                         |
| `make test`         | Runs the test suite (`go test -v ./...`).                                                                     |
| `make lint`         | Runs `golangci-lint`.                                                                                         |
| `make gen`          | Generates mocks via `go generate ./...`.                                                                      |
| `make gen-proto`    | Compiles Protobuf code files using `buf`.                                                                     |
| `make docker-build` | Builds a minimal non-root Docker image.                                                                       |
| `make docker-run`   | Runs the Docker image on port `:8080`.                                                                        |
