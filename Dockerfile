FROM golang:1.11.0
COPY . $GOPATH/src/deathsy/tmga-ga
WORKDIR $GOPATH/src/deathsy/tmga-ga

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

RUN go build -o /go/bin/tmga
COPY ./config.toml /go/bin/config.toml

ENTRYPOINT ["/go/bin/tmga"]

