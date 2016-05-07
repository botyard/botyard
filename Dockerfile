FROM golang:1.6.2

WORKDIR /go/src/github.com/botyard/botyard
ADD . /go/src/github.com/botyard/botyard
RUN go install -v

ENTRYPOINT ["botyard"]
