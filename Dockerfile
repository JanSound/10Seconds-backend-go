FROM golang:1.20-alpine as build

WORKDIR /app

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src/ .

RUN go build -o main .

EXPOSE 8001

CMD ["/app/main"]