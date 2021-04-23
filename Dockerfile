FROM golang:alpine
ENV SRC_DIR=/go/src/github.com/ajarv/go-app

RUN mkdir -p ${SRC_DIR}
ADD . ${SRC_DIR}
WORKDIR ${SRC_DIR}
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh 

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
EXPOSE 8080
EXPOSE 8443
CMD ./main -port ${LISTEN_PORT}  -secure=${LISTEN_SSL}
