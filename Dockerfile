FROM golang:1.7-alpine
ENV SRC_DIR=/go/src/github.com/ajarv/go-web-docker

RUN mkdir -p ${SRC_DIR}
ADD . ${SRC_DIR}
WORKDIR ${SRC_DIR}
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go get -d -v  github.com/go-redis/redis github.com/gorilla/mux gopkg.in/yaml.v2 \
    github.com/thedevsaddam/gojsonq
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest  
RUN apk update && apk upgrade && \
    apk --no-cache add ca-certificates curl && \
    mkdir -p /work 
WORKDIR /work
COPY --from=0 /go/src/github.com/ajarv/go-web-redis/main .
ADD ./static /work/static
ADD ./templates /work/templates
USER 1012
ENV LISTEN_PORT=8080
EXPOSE ${LISTEN_PORT}
CMD ./main --port ${LISTEN_PORT}  