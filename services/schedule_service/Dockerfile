FROM docker.io/golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o schedule-service ./cmd/server

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/schedule-service .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

USER app

EXPOSE 8080

ENTRYPOINT ["./schedule-service"]
