FROM golang:1.21-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .  

RUN GO111MODULE=on go build -o budget-in-api

EXPOSE 8080

CMD ./budget-in-api  
