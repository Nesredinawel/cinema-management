# -------- Builder stage --------
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN go build -o auth-backend .

# -------- Runner stage --------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/auth-backend .

EXPOSE 8081

CMD ["./auth-backend"]
