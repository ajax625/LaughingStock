FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o laughingstock ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/laughingstock .
COPY --from=builder /app/web/dist ./web/dist
EXPOSE 9090
CMD ["./laughingstock"]
