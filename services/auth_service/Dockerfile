FROM docker.io/golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o auth-service ./cmd/server

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations

USER app

EXPOSE 8080

ENTRYPOINT ["./auth-service"]
