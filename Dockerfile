FROM golang:1.7-alpine
USER 1012

ENV SRC_DIR=/go/src/github.com/ajarv/go-app

RUN mkdir -p ${SRC_DIR}
ADD . ${SRC_DIR}
WORKDIR ${SRC_DIR}
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl && \
    go get -d -v  \
    github.com/go-redis/redis \
    github.com/gorilla/mux \
    gopkg.in/yaml.v2 \
    github.com/thedevsaddam/gojsonq

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
EXPOSE 8080
CMD ./main --port 8080  