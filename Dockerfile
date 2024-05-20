FROM golang:1.21-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN mkdir -p templates
RUN go mod download

COPY . .  

RUN GO111MODULE=on go build -o main main.go
FROM alpine:3.19
WORKDIR /app

COPY ./templates/email.html ./templates/
COPY --from=builder /app/main .
COPY dev.env .
COPY prod.env .
COPY --from=0 /app/templates /app/templates

EXPOSE 8080 9090

CMD [ "/app/main" ]