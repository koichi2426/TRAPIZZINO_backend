# ---- Build Stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY src/ src/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src

# ---- Run Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
