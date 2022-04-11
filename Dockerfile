FROM golang:latest

RUN mkdir /agora-assignments

ADD . /agora-assignments

WORKDIR /agora-assignments

RUN go build -o main .

EXPOSE 8000

CMD ["/agora-assignments/main"]