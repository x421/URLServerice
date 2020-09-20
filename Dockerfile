FROM golang:latest
COPY . /go/src/LinksService

WORKDIR /go/src
RUN go get "github.com/go-sql-driver/mysql"
RUN export User=root
RUN export Pass=root
RUN export Ip=127.0.0.1
RUN export PortDB=3306
RUN export PORT=8080

RUN set -ex; CGO_ENABLED=0 GOOS=linux GOARCH=amd64; go build -o ./service LinksService

CMD /bin/bash