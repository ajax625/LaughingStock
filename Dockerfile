FROM node:20-alpine AS web-builder
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.25-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o laughingstock ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=go-builder /app/laughingstock .
COPY --from=web-builder /web/dist ./web/dist
EXPOSE 9090
CMD ["./laughingstock"]
