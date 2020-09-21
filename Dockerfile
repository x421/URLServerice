FROM golang:latest
COPY . /go/src/LinksService

WORKDIR /go/src
RUN go get "github.com/go-sql-driver/mysql"

COPY ./static /go/src/static
RUN set -ex; CGO_ENABLED=1 GOOS=linux GOARCH=amd64; go build -ldflags '-linkmode external -w -extldflags "-static"' -o ./service LinksService

EXPOSE 8080
ENTRYPOINT ./service