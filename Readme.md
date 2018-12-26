## Web App in Go

This is a sample Go App for runtime show case
* Native
* Docker
* OpenShift
* Cloud Foundary




Features - Redis / Bootstrap4 / Go Template but bare bones and none necessary.

## 0. Clone Project

```
git clone https://github.com/ajarv/go-app-docker.git
```

## 1. Run on Local Host

```
cd go-app-docker
run app.go
```

Access the app using curl

```
curl localhost:8080
```

Start Redis instance on localhost

access http://localhost:8080/redis on browser

## 2. Create docker Image

Two phaze docker build

```sh
docker build -t go-app .
```

## 3. Open shift

Assuming you have openshift installed somewhere
And you are logged into openshift with oc.

### 3.1 Deploy on OpenShift plain directly

This won't create redis container

```
oc new-app <this git repo url>
```

### 3.2 Deploy on openshift using template

```
oc process -f  kube-cfg/openshift-app-template.yaml -p APP_NAME=angel  | oc create -f -
```


## 4 Cloud Foundary
