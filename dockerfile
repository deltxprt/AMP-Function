FROM golang:1.19-alpine3.17

RUN mkdir /etc/ampstatus

WORKDIR /etc/ampstatus

COPY ampstatus-azfunction ./ampstatus-azfunction

RUN chmod +x ./ampstatus-azfunction

CMD [ "./ampstatus-azfunction" ]

