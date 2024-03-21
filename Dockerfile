FROM golang:1.21-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .  

RUN GO111MODULE=on go build -o main main.go
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY env_dev.env .

EXPOSE 8080 9090

CMD [ "/app/main" ]