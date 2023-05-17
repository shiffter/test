FROM golang:alpine

WORKDIR /app

COPY ./src /app
COPY src/go.mod src/go.sum ./

RUN go mod tidy

RUN go build -o /main ./cmd/main.go

CMD ["/main"]
