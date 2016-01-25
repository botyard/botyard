FROM golang:1.5.3

WORKDIR /go/src/github.com/botyard/botyard
ADD . /go/src/github.com/botyard/botyard
ENV GOPATH /go/src/github.com/botyard/botyard/Godeps/_workspace:$GOPATH
RUN go install -v

ENTRYPOINT ["botyard"]
