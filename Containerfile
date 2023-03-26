# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

run go build -o amp ./cmd/api

FROM golang:1.20-alpine

WORKDIR /

COPY --from=build /app/amp /amp

EXPOSE 8080

CMD [ "/amp" ]

