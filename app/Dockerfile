FROM golang:1.22.3

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./poster-app

EXPOSE 8080

CMD ["./poster-app"]
