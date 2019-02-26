FROM golang:1.7-alpine
USER 1012
ADD . /tmp/go-app
WORKDIR /tmp/go-app
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go get -d -v  \
    github.com/go-redis/redis \
    github.com/gorilla/mux \
    gopkg.in/yaml.v2 \
    github.com/thedevsaddam/gojsonq \
    github.com/yalp/jsonpath

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest  
ENV SRC_DIR=/go/src/github.com/ajarv/go-app
ENV LISTEN_PORT=8080

RUN apk update && apk upgrade && \
    apk --no-cache add ca-certificates curl && \
    mkdir -p /work 
WORKDIR /work
COPY --from=0 /tmp/go-app/main .
ADD ./static /work/static
ADD ./templates /work/templates
USER 1012
EXPOSE 8080
CMD ./main --port 8080  