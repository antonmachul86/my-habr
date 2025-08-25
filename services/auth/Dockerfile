FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY . .
RUN ls -la
RUN go mod download
RUN go build -o auth-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth-service .
EXPOSE 8001
CMD ["./auth-service"]