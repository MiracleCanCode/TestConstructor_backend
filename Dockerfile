FROM golang:1.23.4-alpine3.21
WORKDIR /app

ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o myapp ./cmd