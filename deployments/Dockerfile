FROM golang:1.23.4-alpine3.20 AS builder
WORKDIR /app

COPY go.mod go.sum ../
RUN go mod download

COPY . .

RUN go build -o /bin/myapp ./cmd

FROM debian:bullseye-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=builder /bin/myapp /bin/myapp
COPY .env.production /app/.env.production

EXPOSE 8080

ENTRYPOINT ["/wait-for-it.sh", "postgres:5432", "--", "/bin/myapp"]
