FROM golang
RUN go get github.com/BurntSushi/toml
RUN go get github.com/samuel/go-zookeeper/zk
WORKDIR /go/src/github.com/dmage/clustr
ADD . .
RUN go get -v ./...
CMD ["/bin/sh", "-c", "/go/bin/clustr"]
