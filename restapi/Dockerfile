FROM golang:latest

ADD . /go/src/restapi

RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2

RUN go install restapi

ENTRYPOINT /go/bin/restapi

EXPOSE 8080
