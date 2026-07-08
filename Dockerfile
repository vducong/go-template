# Build Stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -trimpath -o bin/app ./cmd/app

# Run Stage
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=nonroot:nonroot /app/bin/app /app/app
COPY --from=builder --chown=nonroot:nonroot /app/configs /app/configs

EXPOSE 8080 9090

ENTRYPOINT ["/app/app"]
