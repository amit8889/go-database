FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o go-template ./cmd/go-template
COPY config/config.yaml /app/config/config.yaml
EXPOSE 8001
CMD ["./go-template", "-config", "config/config.yaml"]