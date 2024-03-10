FROM golang:1.21.8

WORKDIR /usr/src/cache

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/cache ./main.go

CMD ["cache"]
