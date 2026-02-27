# ---- Build Stage ----
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/main.go

# ---- Run Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
