# base image
FROM golang:1.14-alpine3.11
ENV CGO_ENABLED 1

RUN apk update && apk add bash inotify-tools git
RUN apk add --update gcc musl-dev

RUN mkdir -p /usr/src/app/backend_go
RUN mkdir /usr/src/app/db
WORKDIR /usr/src/app/backend_go
RUN go get github.com/derekparker/delve/cmd/dlv

COPY ./backend_go ./
RUN go mod vendor
RUN go build -gcflags "all=-N -l" -o /server main.go 

ENTRYPOINT sh start_script.sh

