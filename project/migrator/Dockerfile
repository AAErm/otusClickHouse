ARG GO_IMG=golang:1.22-alpine3.20

FROM $GO_IMG AS builder

ENV GOFLAGS=
WORKDIR /app

COPY . .

RUN go mod download
WORKDIR /app/migrator

RUN GOOS=linux GOARCH=amd64 go build -v -o app ./main.go

CMD ["./app"]
