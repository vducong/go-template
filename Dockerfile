# Step 1: Modules caching
FROM golang:1.18.2-alpine3.16 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.18.2-alpine3.16 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -tags jsoniter -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
COPY --from=builder /app/configs /configs
COPY --from=builder /app/assets /assets
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Step 4: Config timezone
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 5555

ENTRYPOINT ["/app"]
