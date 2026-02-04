# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache make git protobuf ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

# Run Stage
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/configs /app/configs

EXPOSE 8080 9090

ENTRYPOINT ["/app/app"]
