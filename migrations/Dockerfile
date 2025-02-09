FROM golang:1.23.5-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /bin/migrate ./migrations/auto.go

FROM alpine:3.20

RUN apk add --no-cache bash
COPY --from=builder /bin/migrate /bin/migrate
COPY wait-for-it.sh ./wait-for-it.sh

ENTRYPOINT ["/bin/sh", "-c", "./wait-for-it.sh -h postgres -p 5432 -t 60 && /bin/migrate"]
    

    
