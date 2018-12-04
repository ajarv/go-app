# Sample GO app for running in Kubernetes

Features - Redis / Bootstrap4 / Go Template but bare bones and none necessary.

## Clone Project

```
git clone <this git repo url>
```

## Run on Local Host

```
cd go-web-redis-docker
run app.go
```

Access the app using curl

```
curl localhost:8080
```

Start Redis instance on localhost

access http://localhost:8080/redis on browser

## Create docker Image

Two phaze docker build

```sh
docker build -t go-app .
```


# Open shift

Assuming you have openshift installed somewhere 
And you are logged into openshift with oc. 

## Deploy on OpenShift  plain  directly 
This won't create redis container

```
oc new-app <this git repo url>
```

## Deploy on openshift using template

```
oc process -f  openshift-app-template.yaml -p APPLICATION_DOMAIN=go-redis-app.<OPEN SHIFT INFA VIP/IP>.nip.io  | oc create -f -
```
