FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o server cmd/server/server.go

CMD [ "/app/server" ]
