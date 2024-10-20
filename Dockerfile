FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o go-database ./cmd/go-cache
COPY config/config.yaml /app/config/config.yaml
EXPOSE 8001
CMD ["./go-database", "-config", "config/config.yaml"]