FROM golang:1.24

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownling them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o bringten-server main.go

EXPOSE 8080

CMD ["./bringten-server"]
