# Stage 1: Build the binary
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /fabric-aggregator ./cmd/aggregator

# Stage 2: Create a lightweight production image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /fabric-aggregator .

# Run as non-root user for CJIS/SOC2 security best practices
RUN adduser -D fabricuser
USER fabricuser

EXPOSE 8080
CMD ["./fabric-aggregator"]