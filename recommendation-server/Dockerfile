FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o recommendation-server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/recommendation-server .

COPY .env .

EXPOSE 5001

CMD ["./recommendation-server"]
