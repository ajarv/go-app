FROM golang:1.7-alpine
RUN mkdir -p /go/src/github.com/ajarv/go-web-redis/ 
ADD . /go/src/github.com/ajarv/go-web-redis/
WORKDIR /go/src/github.com/ajarv/go-web-redis/
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go get -d -v  github.com/go-redis/redis github.com/gorilla/mux gopkg.in/yaml.v2
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
ENV LISTEN_PORT 8080
EXPOSE ${LISTEN_PORT}
CMD ./main --port ${LISTEN_PORT}  