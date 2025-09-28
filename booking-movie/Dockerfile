# -------- Builder stage --------
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN go build -o booking-movie .

# -------- Runner stage --------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/booking-movie .

EXPOSE 8083

CMD ["./booking-movie"]
