FROM docker.io/golang:1.25.5-alpine AS builder

RUN apk add --no-cache gcc g++ git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o vk-parser ./cmd/parser

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-rus \
    tesseract-ocr-data-eng \
    ca-certificates

WORKDIR /app

COPY --from=builder /app/vk-parser .


USER app

ENTRYPOINT ["./vk-parser"]
