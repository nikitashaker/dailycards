FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY web/templates ./templates
RUN go build -o api ./cmd/main.go
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/api ./api
COPY --from=builder /app/templates ./templates
COPY web ./web
EXPOSE 8080
CMD ["./api"]