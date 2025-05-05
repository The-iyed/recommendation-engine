FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o relation-builder .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/relation-builder .

COPY .env .

EXPOSE 5001

CMD ["./relation-builder"]
