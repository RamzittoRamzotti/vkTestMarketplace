FROM golang:1.23.2-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o /app/marketplace ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/marketplace .
COPY --from=builder /app/.env .
COPY --from=builder /app/storage ./storage
EXPOSE 8080
CMD ["./marketplace"]
