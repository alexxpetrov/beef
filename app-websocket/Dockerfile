# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app-websocket
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /app-websocket/app /app-websocket/cmd/main.go

# Run stage
FROM alpine
WORKDIR /app-websocket
COPY --from=builder /app-websocket/app .

EXPOSE 4000 4600
ENTRYPOINT [ "./app", "-config", "/etc/app-websocket/dev.yaml", "-env", "/etc/app-websocket/.env"]
