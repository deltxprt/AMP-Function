FROM cgr.dev/chainguard/go:latest as build

RUN mkdir /etc/api
RUN mkdir /etc/api/conf

COPY . /etc/api

WORKDIR /etc/api

RUN go build -o amp ./cmd/api

FROM cgr.dev/chainguard/static:latest

COPY --from=build /etc/api/amp /amp

CMD [ "/amp" ]

