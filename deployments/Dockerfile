FROM golang:1.23.5-alpine3.20
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY wait-for-it.sh /app/
RUN chmod +x /app/wait-for-it.sh

EXPOSE 8080:8080
COPY . .
RUN go build -o myapp ./cmd