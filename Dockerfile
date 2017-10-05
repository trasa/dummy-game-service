FROM golang:1.9
MAINTAINER Tony Rasa <t_rasa@wargaming.net>

WORKDIR /go/src/app

COPY . .
RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]

EXPOSE 8080

