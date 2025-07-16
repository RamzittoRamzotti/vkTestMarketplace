FROM golang:1.23.2 AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . /app
RUN go build -o marketplace ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/marketplace .
COPY --from=builder /app/.env .
COPY --from=builder /app/storage ./storage
EXPOSE 8080
CMD ["./marketplace"]