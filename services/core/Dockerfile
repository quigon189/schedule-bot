FROM docker.io/golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o core ./cmd

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/core .

USER app

EXPOSE ${PORT:-8080}

ENTRYPOINT ["./core"]
