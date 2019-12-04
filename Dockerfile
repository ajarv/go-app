FROM golang:1.7-alpine
ENV SRC_DIR=/go/src/github.com/ajarv/go-app

RUN mkdir -p ${SRC_DIR}
ADD . ${SRC_DIR}
WORKDIR ${SRC_DIR}
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go get -d -v  \
    # github.com/go-redis/redis/v7 \
    github.com/gorilla/mux \
    gopkg.in/yaml.v2 \
    github.com/thedevsaddam/gojsonq \
    github.com/yalp/jsonpath

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest  
ENV SRC_DIR=/go/src/github.com/ajarv/go-app
ENV LISTEN_PORT=8080
ENV LISTEN_SSL=f

RUN apk update && apk upgrade && \
    apk --no-cache add ca-certificates curl bash && \
    mkdir -p /work 
WORKDIR /work
COPY --from=0 ${SRC_DIR}/main .
ADD ./static /work/static
ADD ./templates /work/templates
USER 1012
EXPOSE ${LISTEN_PORT}
CMD ./main -port ${LISTEN_PORT}  -secure=${LISTEN_SSL}
