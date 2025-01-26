FROM golang:1.23.5-alpine3.20 AS builder
WORKDIR /app

COPY go.mod go.sum ../
RUN go mod download

COPY . .

RUN go build -o /bin/myapp ./cmd

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache bash curl

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=builder /bin/myapp /bin/myapp

COPY .env.prodaction /app/.env.prodaction

EXPOSE 8080

ENTRYPOINT ["/wait-for-it.sh", "postgres:5432", "--", "/bin/myapp"]
