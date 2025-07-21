FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/marketplace/main.go -o ./docs

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/marketplace

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/main /app/
COPY --from=builder /build/docs /app/docs/

USER appuser

WORKDIR /app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/main", "--health-check"] || exit 1

CMD ["./main"]