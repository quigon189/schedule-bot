FROM docker.io/golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o core ./cmd

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

RUN apk add --no-cache \
    pandoc \
    ca-certificates \
    fontconfig \
    freetype \
    libstdc++ \
    libx11 \
    libxext \
    libxrender \
    libssl3 \
    ttf-freefont \
    ttf-liberation

WORKDIR /app

COPY --from=builder /app/core .


COPY --from=surnet/alpine-wkhtmltopdf:3.23.2-024b2b2-full\
    /bin/wkhtmltopdf \
    /bin/wkhtmltoimage \
    /bin/libwkhtmltox* \
    /bin/

USER app

EXPOSE ${PORT:-8080}

ENTRYPOINT ["./core"]
