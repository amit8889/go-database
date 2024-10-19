FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o go-cache ./cmd/go-cache
COPY config/config.yaml /app/config/config.yaml
EXPOSE 8001
CMD ["./go-cache", "-config", "config/config.yaml"]